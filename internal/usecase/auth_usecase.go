package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/internal/repository"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/pkg/oauth"
	passwordManager "github.com/loak155/techbranch-backend/pkg/password"
	"github.com/loak155/techbranch-backend/pkg/redis"
)

type IAuthUsecase interface {
	Signup(user domain.User) (domain.User, error)
	Signin(email, password string) (accessToken, refreshToken string, expiresIn int, err error)
	Signout(userID int) error
	RefreshToken(refreshToken string) (accessToken string, err error)
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
}

func NewAuthUsecase(repo repository.IUserRepository, jwtAccessTokenManager jwt.JwtManager, jwtRefreshTokenManager jwt.JwtManager, redisAccessTokenManager redis.RedisManager, redisRefreshTokenManager redis.RedisManager, googleManager oauth.GoogleManager) IAuthUsecase {
	return &authUsecase{repo, jwtAccessTokenManager, jwtRefreshTokenManager, redisAccessTokenManager, redisRefreshTokenManager, googleManager}
}

func (usecase *authUsecase) Signup(user domain.User) (domain.User, error) {
	hashedPassword, err := passwordManager.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to hash password: %v", err)
	}
	newUser := domain.User{Username: user.Username, Email: user.Email, Password: hashedPassword}
	if err := usecase.repo.CreateUser(&newUser); err != nil {
		return domain.User{}, err
	}
	return newUser, nil
}

func (usecase *authUsecase) Signin(email, password string) (accessToken, refreshToken string, expiresIn int, err error) {
	user, err := usecase.repo.GetUserByEmail(email)
	if err != nil {
		return "", "", 0, fmt.Errorf("email or password is incorrect")
	}
	if err := passwordManager.CheckPassword(password, user.Password); err != nil {
		return "", "", 0, fmt.Errorf("email or password is incorrect")
	}
	accessToken, accessTokenJti, err := usecase.jwtAccessTokenManager.GenerateToken(int(user.ID))
	if err != nil {
		return "", "", 0, err
	}
	refreshToken, refreshTokenJti, err := usecase.jwtRefreshTokenManager.GenerateToken(int(user.ID))
	if err != nil {
		return "", "", 0, err
	}
	err = usecase.redisAccessTokenManager.Set(context.Background(), strconv.Itoa(int(user.ID)), accessTokenJti)
	if err != nil {
		return "", "", 0, err
	}
	err = usecase.redisRefreshTokenManager.Set(context.Background(), strconv.Itoa(int(user.ID)), refreshTokenJti)
	if err != nil {
		return "", "", 0, err
	}
	expiresIn = usecase.jwtAccessTokenManager.GetExpiresIn()
	return accessToken, refreshToken, expiresIn, nil
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

func (usecase *authUsecase) RefreshToken(refreshToken string) (accessToken string, err error) {
	claims, err := usecase.jwtRefreshTokenManager.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return "", err
	}
	accessToken, accessTokenJti, err := usecase.jwtAccessTokenManager.GenerateToken(userID)
	if err != nil {
		return "", err
	}
	err = usecase.redisAccessTokenManager.Set(context.Background(), claims.Subject, accessTokenJti)
	if err != nil {
		return "", err
	}
	return accessToken, nil
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
		user = &domain.User{Username: userInfo.Name, Email: userInfo.Email}
		if err := usecase.repo.CreateUser(user); err != nil {
			return "", "", 0, err
		}
	}

	accessToken, accessTokenJti, err := usecase.jwtAccessTokenManager.GenerateToken(int(user.ID))
	if err != nil {
		return "", "", 0, err
	}
	refreshToken, refreshTokenJti, err := usecase.jwtRefreshTokenManager.GenerateToken(int(user.ID))
	if err != nil {
		return "", "", 0, err
	}
	err = usecase.redisAccessTokenManager.Set(context.Background(), strconv.Itoa(int(user.ID)), accessTokenJti)
	if err != nil {
		return "", "", 0, err
	}
	err = usecase.redisRefreshTokenManager.Set(context.Background(), strconv.Itoa(int(user.ID)), refreshTokenJti)
	if err != nil {
		return "", "", 0, err
	}
	expiresIn = usecase.jwtAccessTokenManager.GetExpiresIn()
	return accessToken, refreshToken, expiresIn, nil
}
