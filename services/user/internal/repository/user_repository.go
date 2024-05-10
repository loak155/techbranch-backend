package repository

import (
	"github.com/loak155/techbranch-backend/services/user/internal/domain"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *domain.User) error
	GetUser(id int) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	ListUsers() (*[]domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (repo *userRepository) CreateUser(user *domain.User) error {
	err := repo.db.Create(user).Error
	return err
}

func (repo *userRepository) GetUser(id int) (*domain.User, error) {
	user := &domain.User{}
	err := repo.db.First(user, id).Error
	return user, err
}

func (repo *userRepository) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	err := repo.db.Where("email=?", email).First(user).Error
	return user, err
}

func (repo *userRepository) ListUsers() (*[]domain.User, error) {
	users := &[]domain.User{}
	err := repo.db.Find(&users).Error
	return users, err
}

func (repo *userRepository) UpdateUser(user *domain.User) error {
	err := repo.db.Save(user).Error
	return err
}

func (repo *userRepository) DeleteUser(id int) error {
	err := repo.db.Delete(&domain.User{}, id).Error
	return err
}
