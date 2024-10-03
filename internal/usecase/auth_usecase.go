package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/internal/repository"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/pkg/mail"
	"github.com/loak155/techbranch-backend/pkg/oauth"
	passwordManager "github.com/loak155/techbranch-backend/pkg/password"
	"github.com/loak155/techbranch-backend/pkg/redis"
	"github.com/loak155/techbranch-backend/pkg/uuid"
)

type IAuthUsecase interface {
	PreSignup(user domain.User) error
	Signup(token string) error
	Signin(email, password string) (accessToken, refreshToken string, accessTokenExpiresIn, refreshTokenExpiresIn int, err error)
	Signout(userID int) error
	RefreshToken(refreshToken string) (accessToken string, accessTokenExpiresIn int, err error)
	GetSigninUser(userID int) (domain.User, error)
	GetGoogleLoginURL() string
	GoogleLoginCallback(state, code string) (accessToken, refreshToken string, expiresIn int, err error)
}

type authUsecase struct {
	repo                     repository.IUserRepository
	jwtAccessTokenManager    jwt.JwtManager
	jwtRefreshTokenManager   jwt.JwtManager
	redisAccessTokenManager  redis.RedisManager
	redisRefreshTokenManager redis.RedisManager
	googleManager            oauth.GoogleManager
	presignupRedisManager    redis.RedisManager
	presignupMailManager     mail.PresignupMailManager
}

func NewAuthUsecase(repo repository.IUserRepository, jwtAccessTokenManager jwt.JwtManager, jwtRefreshTokenManager jwt.JwtManager, redisAccessTokenManager redis.RedisManager, redisRefreshTokenManager redis.RedisManager, googleManager oauth.GoogleManager, presignupRedisManager redis.RedisManager, presignupMailManager mail.PresignupMailManager) IAuthUsecase {
	return &authUsecase{repo, jwtAccessTokenManager, jwtRefreshTokenManager, redisAccessTokenManager, redisRefreshTokenManager, googleManager, presignupRedisManager, presignupMailManager}
}

func (usecase *authUsecase) PreSignup(user domain.User) error {
	u, err := usecase.repo.GetUserByEmail(user.Email)
	if err == nil {
		if u.Email == user.Email {
			return fmt.Errorf("user already exists")
		} else {
			return fmt.Errorf("failed to get user: %v", err)
		}
	}

	hashedPassword, err := passwordManager.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	b, err := json.Marshal(domain.User{Username: user.Username, Email: user.Email, Password: hashedPassword})
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}
	userString := string(b)

	uuid := uuid.NewUUID()
	if err := usecase.presignupRedisManager.Set(context.Background(), uuid, userString); err != nil {
		return fmt.Errorf("failed to set redis: %v", err)
	}

	if err := usecase.presignupMailManager.SendPreSignUpMail([]string{user.Email}, user.Username, uuid); err != nil {
		return fmt.Errorf("failed to send mail: %v", err)
	}

	return nil
}

func (usecase *authUsecase) Signup(token string) error {
	userString, err := usecase.presignupRedisManager.Get(context.Background(), token)
	if err != nil {
		return fmt.Errorf("failed to get redis: %v", err)
	}

	var user domain.User
	if err := json.Unmarshal([]byte(userString), &user); err != nil {
		return fmt.Errorf("failed to unmarshal user: %v", err)
	}

	if err := usecase.repo.CreateUser(&user); err != nil {
		return err
	}
	return nil
}

func (usecase *authUsecase) Signin(email, password string) (accessToken, refreshToken string, accessTokenExpiresIn, refreshTokenExpiresIn int, err error) {
	user, err := usecase.repo.GetUserByEmail(email)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("email or password is incorrect")
	}
	if err := passwordManager.CheckPassword(password, user.Password); err != nil {
		return "", "", 0, 0, fmt.Errorf("email or password is incorrect")
	}
	accessToken, accessTokenJti, err := usecase.jwtAccessTokenManager.GenerateToken(int(user.ID))
	if err != nil {
		return "", "", 0, 0, err
	}
	refreshToken, refreshTokenJti, err := usecase.jwtRefreshTokenManager.GenerateToken(int(user.ID))
	if err != nil {
		return "", "", 0, 0, err
	}
	err = usecase.redisAccessTokenManager.Set(context.Background(), strconv.Itoa(int(user.ID)), accessTokenJti)
	if err != nil {
		return "", "", 0, 0, err
	}
	err = usecase.redisRefreshTokenManager.Set(context.Background(), strconv.Itoa(int(user.ID)), refreshTokenJti)
	if err != nil {
		return "", "", 0, 0, err
	}
	accessTokenExpiresIn = usecase.jwtAccessTokenManager.GetExpiresIn()
	refreshTokenExpiresIn = usecase.jwtRefreshTokenManager.GetExpiresIn()
	return accessToken, refreshToken, accessTokenExpiresIn, refreshTokenExpiresIn, nil
}

func (usecase *authUsecase) Signout(userID int) error {
	if err := usecase.redisAccessTokenManager.Del(context.Background(), strconv.Itoa(userID)); err != nil {
		return err
	}
	if err := usecase.redisRefreshTokenManager.Del(context.Background(), strconv.Itoa(userID)); err != nil {
		return err
	}
	return nil
}

func (usecase *authUsecase) RefreshToken(refreshToken string) (accessToken string, accessTokenExpiresIn int, err error) {
	claims, err := usecase.jwtRefreshTokenManager.ValidateToken(refreshToken)
	if err != nil {
		return "", 0, err
	}
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return "", 0, err
	}
	accessToken, accessTokenJti, err := usecase.jwtAccessTokenManager.GenerateToken(userID)
	if err != nil {
		return "", 0, err
	}
	err = usecase.redisAccessTokenManager.Set(context.Background(), claims.Subject, accessTokenJti)
	if err != nil {
		return "", 0, err
	}

	accessTokenExpiresIn = usecase.jwtAccessTokenManager.GetExpiresIn()
	return accessToken, accessTokenExpiresIn, nil
}

func (usecase *authUsecase) GetSigninUser(userID int) (domain.User, error) {
	user, err := usecase.repo.GetUser(userID)
	if err != nil {
		return domain.User{}, err
	}
	return *user, nil
}

func (usecase *authUsecase) GetGoogleLoginURL() string {
	return usecase.googleManager.GetLoginURL()
}

func (usecase *authUsecase) GoogleLoginCallback(state, code string) (accessToken, refreshToken string, expiresIn int, err error) {
	if !usecase.googleManager.CheckState(state) {
		return "", "", 0, fmt.Errorf("invalid state")
	}
	token, err := usecase.googleManager.GetAccessToken(code)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to get access token: %v", err)
	}
	userInfo, err := usecase.googleManager.GetUserInfo(token)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to get user info: %v", err)
	}

	// check if user exists
	user, err := usecase.repo.GetUserByEmail(userInfo.Email)
	// if not exist, create user
	if err != nil {
		user = &domain.User{Username: userInfo.Name, Email: userInfo.Email, GoogleID: userInfo.ID}
		if err := usecase.repo.CreateUser(user); err != nil {
			return "", "", 0, fmt.Errorf("failed to create user: %v", err)
		}
		// if exist and google id is empty, update google id
	} else if user.GoogleID == "" {
		user.GoogleID = userInfo.ID
		if err := usecase.repo.UpdateUser(user); err != nil {
			return "", "", 0, fmt.Errorf("failed to update user: %v", err)
		}
	}

	accessToken, accessTokenJti, err := usecase.jwtAccessTokenManager.GenerateToken(int(user.ID))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate access token: %v", err)
	}
	refreshToken, refreshTokenJti, err := usecase.jwtRefreshTokenManager.GenerateToken(int(user.ID))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate refresh token: %v", err)
	}
	err = usecase.redisAccessTokenManager.Set(context.Background(), strconv.Itoa(int(user.ID)), accessTokenJti)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to set access token: %v", err)
	}
	err = usecase.redisRefreshTokenManager.Set(context.Background(), strconv.Itoa(int(user.ID)), refreshTokenJti)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to set refresh token: %v", err)
	}
	expiresIn = usecase.jwtAccessTokenManager.GetExpiresIn()
	return accessToken, refreshToken, expiresIn, nil
}
