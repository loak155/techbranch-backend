package usecase

import (
	"context"

	"github.com/loak155/techbranch-backend/pkg/password"
	"github.com/loak155/techbranch-backend/services/user/internal/domain"
	"github.com/loak155/techbranch-backend/services/user/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IUserUsecase interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (bool, error)
	DeleteUser(ctx context.Context, id int) (bool, error)
}

type userUsecase struct {
	repo repository.IUserRepository
}

func NewUserUsecase(repo repository.IUserRepository) IUserUsecase {
	return &userUsecase{repo}
}

func (usecase *userUsecase) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	hashedPassword, err := password.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}
	newUser := domain.User{Username: user.Username, Email: user.Email, Password: hashedPassword}
	if err := usecase.repo.CreateUser(&newUser); err != nil {
		return domain.User{}, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	return newUser, nil
}

func (usecase *userUsecase) GetUser(ctx context.Context, id int) (domain.User, error) {
	user, err := usecase.repo.GetUser(id)
	if err != nil {
		return domain.User{}, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	return *user, nil
}

func (usecase *userUsecase) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := usecase.repo.GetUserByEmail(email)
	if err != nil {
		return domain.User{}, status.Errorf(codes.Internal, "failed to get user by email: %v", err)
	}
	return *user, nil
}

func (usecase *userUsecase) ListUsers(ctx context.Context) ([]domain.User, error) {
	users, err := usecase.repo.ListUsers()
	if err != nil {
		return []domain.User{}, status.Errorf(codes.Internal, "failed to get user list: %v", err)
	}
	return *users, nil
}

func (usecase *userUsecase) UpdateUser(ctx context.Context, user domain.User) (bool, error) {
	hashedPassword, err := password.HashPassword(user.Password)
	if err != nil {
		return false, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}
	updatedUser := domain.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
	}
	if err := usecase.repo.UpdateUser(&updatedUser); err != nil {
		return false, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}
	return true, nil
}

func (usecase *userUsecase) DeleteUser(ctx context.Context, id int) (bool, error) {
	if err := usecase.repo.DeleteUser(id); err != nil {
		return false, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	return true, nil
}
