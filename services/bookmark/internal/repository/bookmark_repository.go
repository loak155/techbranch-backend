package repository

import (
	"github.com/loak155/techbranch-backend/services/bookmark/internal/domain"
	"gorm.io/gorm"
)

type IBookmarkRepository interface {
	CreateBookmark(bookmark *domain.Bookmark) error
	UpdateBookmarkWithUnscoped(id int) error
	GetBookmark(id int) (*domain.Bookmark, error)
	GetBookmarkByUserIDAndArticleID(userID, articleID int) (*domain.Bookmark, error)
	GetBookmarkByUserIDAndArticleIDWithUnscoped(userID, articleID int) (*domain.Bookmark, error)
	ListBookmarks() (*[]domain.Bookmark, error)
	ListBookmarksByUserID(userID int) (*[]domain.Bookmark, error)
	ListBookmarksByArticleID(articleID int) (*[]domain.Bookmark, error)
	DeleteBookmark(id int) error
	DeleteBookmarkByUserIDAndArticleID(userID, articleID int) error
	DeleteBookmarkByUserID(UserID int) error
	UpdateBookmarkByUserIDWithUnscoped(UserID int) error
	DeleteBookmarkByArticleID(ArticleID int) error
	UpdateBookmarkByArticleIDWithUnscoped(ArticleID int) error
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

func (repo *bookmarkRepository) UpdateBookmarkWithUnscoped(id int) error {
	err := repo.db.Unscoped().Model(&domain.Bookmark{}).Where(id).Update("deleted_at", nil).Error
	return err
}

func (repo *bookmarkRepository) GetBookmark(id int) (*domain.Bookmark, error) {
	bookmark := &domain.Bookmark{}
	err := repo.db.First(bookmark, id).Error
	return bookmark, err
}

func (repo *bookmarkRepository) GetBookmarkByUserIDAndArticleID(userID, articleID int) (*domain.Bookmark, error) {
	bookmark := &domain.Bookmark{}
	err := repo.db.Where("user_id=? AND article_id=?", userID, articleID).First(bookmark).Error
	return bookmark, err
}

func (repo *bookmarkRepository) GetBookmarkByUserIDAndArticleIDWithUnscoped(userID, articleID int) (*domain.Bookmark, error) {
	bookmark := &domain.Bookmark{}
	err := repo.db.Unscoped().Where("user_id=? AND article_id=?", userID, articleID).First(bookmark).Error
	return bookmark, err
}

func (repo *bookmarkRepository) ListBookmarks() (*[]domain.Bookmark, error) {
	bookmarks := &[]domain.Bookmark{}
	err := repo.db.Find(bookmarks).Error
	return bookmarks, err
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

func (repo *bookmarkRepository) DeleteBookmark(id int) error {
	err := repo.db.Delete(&domain.Bookmark{}, id).Error
	return err
}

func (repo *bookmarkRepository) DeleteBookmarkByUserIDAndArticleID(userID, articleID int) error {
	err := repo.db.Where("user_id=? AND article_id=?", userID, articleID).Delete(&domain.Bookmark{}).Error
	return err
}

func (repo *bookmarkRepository) DeleteBookmarkByUserID(userID int) error {
	err := repo.db.Delete(&domain.Bookmark{}, "user_id=?", userID).Error
	return err
}

func (repo *bookmarkRepository) UpdateBookmarkByUserIDWithUnscoped(userID int) error {
	err := repo.db.Unscoped().Model(&domain.Bookmark{}).Where("user_id", userID).Update("deleted_at", nil).Error
	return err
}

func (repo *bookmarkRepository) DeleteBookmarkByArticleID(articleID int) error {
	err := repo.db.Delete(&domain.Bookmark{}, "article_id=?", articleID).Error
	return err
}

func (repo *bookmarkRepository) UpdateBookmarkByArticleIDWithUnscoped(articleID int) error {
	err := repo.db.Unscoped().Model(&domain.Bookmark{}).Where("article_id", articleID).Update("deleted_at", nil).Error
	return err
}
