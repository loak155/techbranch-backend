package usecase

import (
	"context"

	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/pkg/oauth"
	pb "github.com/loak155/techbranch-backend/services/user/proto"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IAuthUsecase interface {
	Signup(ctx context.Context, username string, email string, password string) (int, error)
	Signin(ctx context.Context, email string, password string) (string, error)
	GetSigninUser(ctx context.Context, userId int) (username string, email string, error error)
	GenerateToken(ctx context.Context, userId int) (string, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
	RefreshToken(ctx context.Context, token string) (string, error)
	GetGoogleLoginURL(ctx context.Context) (string, error)
	GoogleLoginCallback(ctx context.Context, state string, code string) (string, error)
}

type authUsecase struct {
	userClient pb.UserServiceClient
	jwtManager jwt.JwtManager
	google     oauth.Google
}

func NewAuthUsecase(userClient pb.UserServiceClient, jwtManager jwt.JwtManager, google oauth.Google) IAuthUsecase {
	return &authUsecase{userClient, jwtManager, google}
}

func (usecase *authUsecase) Signup(ctx context.Context, username string, email string, password string) (int, error) {
	req := pb.CreateUserRequest{
		Username: username,
		Email:    email,
		Password: password,
	}
	res, err := usecase.userClient.CreateUser(ctx, &req)
	if err != nil {
		return 0, err
	}
	return int(res.User.Id), nil
}

func (usecase *authUsecase) Signin(ctx context.Context, email string, password string) (string, error) {
	req := pb.GetUserByEmailRequest{Email: email}
	res, err := usecase.userClient.GetUserByEmail(ctx, &req)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(res.User.Password), []byte(password))
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "unmatched password: %v", err)
	}
	token, err := usecase.jwtManager.Generate(int(res.User.Id))
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}
	return token, nil
}

func (usecase *authUsecase) GetSigninUser(ctx context.Context, userId int) (username string, email string, err error) {
	res, err := usecase.userClient.GetUser(ctx, &pb.GetUserRequest{Id: int32(userId)})
	if err != nil {
		return "", "", err
	}
	return res.User.Username, res.User.Email, err
}

func (usecase *authUsecase) GenerateToken(ctx context.Context, userId int) (string, error) {
	token, err := usecase.jwtManager.Generate(userId)
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}
	return token, nil
}

func (usecase *authUsecase) ValidateToken(ctx context.Context, token string) (bool, error) {
	_, err := usecase.jwtManager.ValidateToken(token)
	if err != nil {
		return false, status.Errorf(codes.Internal, "invalid token: %v", err)
	}
	return true, nil
}

func (usecase *authUsecase) RefreshToken(ctx context.Context, token string) (string, error) {
	claims, err := usecase.jwtManager.ValidateToken(token)
	if err != nil {
		return "", status.Errorf(codes.Internal, "invalid token: %v", err)
	}
	refreshToken, err := usecase.jwtManager.Generate(claims.UserId)
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}
	return refreshToken, nil
}

func (usecase *authUsecase) GetGoogleLoginURL(ctx context.Context) (string, error) {
	return usecase.google.GetLoginURL(), nil
}

func (usecase *authUsecase) GoogleLoginCallback(ctx context.Context, state string, code string) (string, error) {
	if !usecase.google.CheckState(state) {
		return "", status.Errorf(codes.InvalidArgument, "invalid state")
	}
	token, err := usecase.google.GetAccessToken(code)
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to get access token: %v", err)
	}
	userInfo, err := usecase.google.GetUserInfo(token)
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to get user info: %v", err)
	}

	userId := 0
	// check if user exists
	resGetUser, err := usecase.userClient.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{Email: userInfo.Email})
	// if not exist, create user
	if err != nil {
		resCreateUser, err := usecase.userClient.CreateUser(ctx, &pb.CreateUserRequest{
			Username: userInfo.Name,
			Email:    userInfo.Email,
			Password: "",
		})
		if err != nil {
			return "", err
		}
		userId = int(resCreateUser.User.Id)
	} else {
		userId = int(resGetUser.User.Id)
	}

	jwtToken, err := usecase.jwtManager.Generate(userId)
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}
	return jwtToken, nil
}
