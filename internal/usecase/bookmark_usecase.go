package usecase

import (
	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/internal/repository"
)

type IBookmarkUsecase interface {
	CreateBookmark(bookmark domain.Bookmark) (domain.Bookmark, error)
	GetBookmarkCountByArticleID(articleID int) (int, error)
	ListBookmarksByUserID(userID int) ([]domain.Bookmark, error)
	ListBookmarksByArticleID(articleID int) ([]domain.Bookmark, error)
	DeleteBookmarkByUserIDAndArticleID(userID, articleID int) error
	DeleteBookmarkByUserID(userID int) error
	DeleteBookmarkByArticleID(articleID int) error
}

type bookmarkUsecase struct {
	repo repository.IBookmarkRepository
}

func NewBookmarkUsecase(repo repository.IBookmarkRepository) IBookmarkUsecase {
	return &bookmarkUsecase{repo}
}

func (usecase *bookmarkUsecase) CreateBookmark(bookmark domain.Bookmark) (domain.Bookmark, error) {
	if err := usecase.repo.CreateBookmark(&bookmark); err != nil {
		return domain.Bookmark{}, err
	}
	return bookmark, nil
}

func (usecase *bookmarkUsecase) GetBookmarkCountByArticleID(articleID int) (int, error) {
	count, err := usecase.repo.GetBookmarkCountByArticleID(articleID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (usecase *bookmarkUsecase) ListBookmarksByUserID(userID int) ([]domain.Bookmark, error) {
	bookmarks, err := usecase.repo.ListBookmarksByUserID(userID)
	if err != nil {
		return []domain.Bookmark{}, err
	}
	return *bookmarks, nil
}

func (usecase *bookmarkUsecase) ListBookmarksByArticleID(articleID int) ([]domain.Bookmark, error) {
	bookmarks, err := usecase.repo.ListBookmarksByArticleID(articleID)
	if err != nil {
		return []domain.Bookmark{}, err
	}
	return *bookmarks, nil
}

func (usecase *bookmarkUsecase) DeleteBookmarkByUserIDAndArticleID(userID, articleID int) error {
	err := usecase.repo.DeleteBookmarkByUserIDAndArticleID(userID, articleID)
	return err
}

func (usecase *bookmarkUsecase) DeleteBookmarkByUserID(userID int) error {
	err := usecase.repo.DeleteBookmarkByUserID(userID)
	return err
}

func (usecase *bookmarkUsecase) DeleteBookmarkByArticleID(articleID int) error {
	err := usecase.repo.DeleteBookmarkByArticleID(articleID)
	return err
}
