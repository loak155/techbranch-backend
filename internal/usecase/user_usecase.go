package usecase

import (
	"fmt"

	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/internal/repository"
	"github.com/loak155/techbranch-backend/pkg/password"
)

type IUserUsecase interface {
	CreateUser(user domain.User) (domain.User, error)
	GetUser(id int) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	ListUsers(offset, limit int) ([]domain.User, error)
	UpdateUser(user domain.User) (domain.User, error)
	DeleteUser(id int) error
}

type userUsecase struct {
	repo repository.IUserRepository
}

func NewUserUsecase(repo repository.IUserRepository) IUserUsecase {
	return &userUsecase{repo}
}

func (usecase *userUsecase) CreateUser(user domain.User) (domain.User, error) {
	hashedPassword, err := password.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to hash password: %v", err)
	}
	newUser := domain.User{Username: user.Username, Email: user.Email, Password: hashedPassword}
	if err := usecase.repo.CreateUser(&newUser); err != nil {
		return domain.User{}, err
	}
	return newUser, nil
}

func (usecase *userUsecase) GetUser(id int) (domain.User, error) {
	user, err := usecase.repo.GetUser(id)
	if err != nil {
		return domain.User{}, err
	}
	return *user, nil
}

func (usecase *userUsecase) GetUserByEmail(email string) (domain.User, error) {
	user, err := usecase.repo.GetUserByEmail(email)
	if err != nil {
		return domain.User{}, err
	}
	return *user, nil
}

func (usecase *userUsecase) ListUsers(offset, limit int) ([]domain.User, error) {
	users, err := usecase.repo.ListUsers(offset, limit)
	if err != nil {
		return []domain.User{}, err
	}
	return *users, nil
}

func (usecase *userUsecase) UpdateUser(user domain.User) (domain.User, error) {
	updatedUser := domain.User{}
	if user.Password == "" {
		updatedUser = domain.User{ID: user.ID, Username: user.Username, Email: user.Email}
	} else {
		hashedPassword, err := password.HashPassword(user.Password)
		if err != nil {
			return domain.User{}, fmt.Errorf("failed to hash password: %v", err)
		}
		updatedUser = domain.User{ID: user.ID, Username: user.Username, Email: user.Email, Password: hashedPassword}
	}
	if err := usecase.repo.UpdateUser(&updatedUser); err != nil {
		return domain.User{}, err
	}
	return updatedUser, nil
}

func (usecase *userUsecase) DeleteUser(id int) error {
	err := usecase.repo.DeleteUser(id)
	return err
}
