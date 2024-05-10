package repository

import (
	"github.com/loak155/techbranch-backend/services/comment/internal/domain"
	"gorm.io/gorm"
)

type ICommentRepository interface {
	CreateComment(comment *domain.Comment) error
	GetComment(id int) (*domain.Comment, error)
	ListCommentsByArticleID(articleID int) (*[]domain.Comment, error)
	UpdateComment(comment *domain.Comment) error
	DeleteComment(id int) error
	DeleteCommentByUserID(userID int) error
	UpdateByUserIDWithUnscoped(userID int) error
	DeleteCommentByArticleID(articleID int) error
	UpdateByArticleIDWithUnscoped(articleID int) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) ICommentRepository {
	return &commentRepository{db}
}

func (repo *commentRepository) CreateComment(comment *domain.Comment) error {
	err := repo.db.Create(comment).Error
	return err
}

func (repo *commentRepository) GetComment(id int) (*domain.Comment, error) {
	comment := &domain.Comment{}
	err := repo.db.First(comment, id).Error
	return comment, err
}

func (repo *commentRepository) ListCommentsByArticleID(articleID int) (*[]domain.Comment, error) {
	comments := &[]domain.Comment{}
	err := repo.db.Where("article_id=?", articleID).Find(comments).Error
	return comments, err
}

func (repo *commentRepository) UpdateComment(comment *domain.Comment) error {
	err := repo.db.Save(comment).Error
	return err
}

func (repo *commentRepository) DeleteComment(id int) error {
	err := repo.db.Delete(&domain.Comment{}, id).Error
	return err
}

func (repo *commentRepository) DeleteCommentByUserID(userID int) error {
	err := repo.db.Delete(&domain.Comment{}, "user_id=?", userID).Error
	return err
}

func (repo *commentRepository) UpdateByUserIDWithUnscoped(userID int) error {
	err := repo.db.Unscoped().Model(&domain.Comment{}).Where("user_id", userID).Update("deleted_at", nil).Error
	return err
}

func (repo *commentRepository) DeleteCommentByArticleID(articleID int) error {
	err := repo.db.Delete(&domain.Comment{}, "article_id=?", articleID).Error
	return err
}

func (repo *commentRepository) UpdateByArticleIDWithUnscoped(articleID int) error {
	err := repo.db.Unscoped().Model(&domain.Comment{}).Where("article_id", articleID).Update("deleted_at", nil).Error
	return err
}
