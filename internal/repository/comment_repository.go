package repository

import (
	"github.com/loak155/techbranch-backend/internal/domain"
	"gorm.io/gorm"
)

type ICommentRepository interface {
	CreateComment(comment *domain.Comment) error
	ListCommentsByUserID(userID int) (*[]domain.Comment, error)
	ListCommentsByArticleID(articleID int) (*[]domain.Comment, error)
	DeleteComment(id int) error
	DeleteCommentByUserIDAndArticleID(userID, articleID int) error
	DeleteCommentByUserID(UserID int) error
	DeleteCommentByArticleID(ArticleID int) error
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

func (repo *commentRepository) ListCommentsByUserID(userID int) (*[]domain.Comment, error) {
	comments := &[]domain.Comment{}
	err := repo.db.Where("user_id=?", userID).Find(comments).Error
	return comments, err
}

func (repo *commentRepository) ListCommentsByArticleID(articleID int) (*[]domain.Comment, error) {
	comments := &[]domain.Comment{}
	err := repo.db.Where("article_id=?", articleID).Find(comments).Error
	return comments, err
}

func (repo *commentRepository) DeleteComment(id int) error {
	err := repo.db.Delete(&domain.Comment{}, id).Error
	return err
}

func (repo *commentRepository) DeleteCommentByUserIDAndArticleID(userID, articleID int) error {
	err := repo.db.Where("user_id=? AND article_id=?", userID, articleID).Delete(&domain.Comment{}).Error
	return err
}

func (repo *commentRepository) DeleteCommentByUserID(userID int) error {
	err := repo.db.Delete(&domain.Comment{}, "user_id=?", userID).Error
	return err
}

func (repo *commentRepository) DeleteCommentByArticleID(articleID int) error {
	err := repo.db.Delete(&domain.Comment{}, "article_id=?", articleID).Error
	return err
}
