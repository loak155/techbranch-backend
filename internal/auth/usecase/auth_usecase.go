package usecase

import (
	"context"

	"github.com/loak155/techbranch-backend/internal/pkg/jwt"
	pb "github.com/loak155/techbranch-backend/proto/user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IAuthUsecase interface {
	Signup(ctx context.Context, username string, email string, password string) (int, error)
	Signin(ctx context.Context, email string, password string) (string, error)
	GenerateToken(ctx context.Context, user_id int) (string, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
	RefreshToken(ctx context.Context, token string) (string, error)
}

type authUsecase struct {
	userClient pb.UserServiceClient
	jwtManager jwt.JwtManager
}

func NewAuthUsecase(userClient pb.UserServiceClient, jwtManager jwt.JwtManager) IAuthUsecase {
	return &authUsecase{userClient, jwtManager}
}

func (usecase *authUsecase) Signup(ctx context.Context, username string, email string, password string) (int, error) {
	req := pb.CreateUserRequest{
		User: &pb.User{
			Username: username,
			Email:    email,
			Password: password,
		},
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

func (usecase *authUsecase) GenerateToken(ctx context.Context, user_id int) (string, error) {
	token, err := usecase.jwtManager.Generate(user_id)
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
