package repository

import (
	"github.com/loak155/techbranch-backend/internal/domain"
	"gorm.io/gorm"
)

type IBookmarkRepository interface {
	CreateBookmark(bookmark *domain.Bookmark) error
	GetBookmarkCountByArticleID(articleID int) (int, error)
	ListBookmarksByUserID(userID int) (*[]domain.Bookmark, error)
	ListBookmarksByArticleID(articleID int) (*[]domain.Bookmark, error)
	DeleteBookmarkByUserIDAndArticleID(userID, articleID int) error
	DeleteBookmarkByUserID(UserID int) error
	DeleteBookmarkByArticleID(ArticleID int) error
}

type bookmarkRepository struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) IBookmarkRepository {
	return &bookmarkRepository{db}
}

func (repo *bookmarkRepository) CreateBookmark(bookmark *domain.Bookmark) error {
	err := repo.db.Create(bookmark).Error
	return err
}

func (repo *bookmarkRepository) GetBookmarkCountByArticleID(articleID int) (int, error) {
	var count int64
	err := repo.db.Model(&domain.Bookmark{}).Where("article_id=?", articleID).Count(&count).Error
	return int(count), err
}

func (repo *bookmarkRepository) ListBookmarksByUserID(userID int) (*[]domain.Bookmark, error) {
	bookmarks := &[]domain.Bookmark{}
	err := repo.db.Where("user_id=?", userID).Find(bookmarks).Error
	return bookmarks, err
}

func (repo *bookmarkRepository) ListBookmarksByArticleID(articleID int) (*[]domain.Bookmark, error) {
	bookmarks := &[]domain.Bookmark{}
	err := repo.db.Where("article_id=?", articleID).Find(bookmarks).Error
	return bookmarks, err
}

func (repo *bookmarkRepository) DeleteBookmarkByUserIDAndArticleID(userID, articleID int) error {
	err := repo.db.Where("user_id=? AND article_id=?", userID, articleID).Delete(&domain.Bookmark{}).Error
	return err
}

func (repo *bookmarkRepository) DeleteBookmarkByUserID(userID int) error {
	err := repo.db.Delete(&domain.Bookmark{}, "user_id=?", userID).Error
	return err
}

func (repo *bookmarkRepository) DeleteBookmarkByArticleID(articleID int) error {
	err := repo.db.Delete(&domain.Bookmark{}, "article_id=?", articleID).Error
	return err
}
