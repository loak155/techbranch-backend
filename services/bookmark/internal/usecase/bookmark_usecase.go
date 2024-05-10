package usecase

import (
	"context"

	pb "github.com/loak155/techbranch-backend/services/article/proto"
	"github.com/loak155/techbranch-backend/services/bookmark/internal/domain"
	"github.com/loak155/techbranch-backend/services/bookmark/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IBookmarkUsecase interface {
	CreateBookmark(ctx context.Context, bookmark domain.Bookmark) (domain.Bookmark, error)
	GetBookmark(ctx context.Context, id int) (domain.Bookmark, error)
	GetBookmarkCountByArticleID(ctx context.Context, articleID int) (int, error)
	ListBookmarks(ctx context.Context) ([]domain.Bookmark, error)
	ListBookmarksByUserID(ctx context.Context, UserID int) ([]domain.Bookmark, error)
	ListBookmarksByArticleID(ctx context.Context, articleID int) ([]domain.Bookmark, error)
	DeleteBookmark(ctx context.Context, bookmark domain.Bookmark) (bool, error)
	DeleteBookmarkByUserID(ctx context.Context, UserID int) (bool, error)
	DeleteBookmarkByUserIDCompensate(ctx context.Context, UserID int) (bool, error)
	DeleteBookmarkByArticleID(ctx context.Context, ArticleID int) (bool, error)
	DeleteBookmarkByArticleIDCompensate(ctx context.Context, ArticleID int) (bool, error)
}

type bookmarkUsecase struct {
	repo          repository.IBookmarkRepository
	articleClient pb.ArticleServiceClient
}

func NewBookmarkUsecase(repo repository.IBookmarkRepository, articleClient pb.ArticleServiceClient) IBookmarkUsecase {
	return &bookmarkUsecase{repo, articleClient}
}

func (usecase *bookmarkUsecase) CreateBookmark(ctx context.Context, bookmark domain.Bookmark) (domain.Bookmark, error) {
	storeBookmark, err := usecase.repo.GetBookmarkByUserIDAndArticleIDWithUnscoped(int(bookmark.UserID), int(bookmark.ArticleID))
	if storeBookmark.ID != 0 && err == nil {
		// if bookmark already exists, update it
		if err := usecase.repo.UpdateBookmarkWithUnscoped(int(storeBookmark.ID)); err != nil {
			return domain.Bookmark{}, status.Errorf(codes.Internal, "failed to compensate bookmark deletion: %v", err)
		}
		return *storeBookmark, nil
	} else {
		// if bookmark does not exist, create it
		if err := usecase.repo.CreateBookmark(&bookmark); err != nil {
			return domain.Bookmark{}, status.Errorf(codes.Internal, "failed to create bookmark: %v", err)
		}
		return bookmark, nil
	}
}

func (usecase *bookmarkUsecase) GetBookmark(ctx context.Context, id int) (domain.Bookmark, error) {
	bookmark, err := usecase.repo.GetBookmark(id)
	if err != nil {
		return domain.Bookmark{}, status.Errorf(codes.Internal, "failed to get bookmark: %v", err)
	}
	return *bookmark, nil
}

func (usecase *bookmarkUsecase) GetBookmarkCountByArticleID(ctx context.Context, articleID int) (int, error) {
	bookmarks, err := usecase.repo.ListBookmarksByArticleID(articleID)
	if err != nil {
		return 0, status.Errorf(codes.Internal, "failed to get bookmark list: %v", err)
	}
	return len(*bookmarks), nil
}

func (usecase *bookmarkUsecase) ListBookmarks(ctx context.Context) ([]domain.Bookmark, error) {
	bookmarks, err := usecase.repo.ListBookmarks()
	if err != nil {
		return []domain.Bookmark{}, status.Errorf(codes.Internal, "failed to get bookmark list: %v", err)
	}
	return *bookmarks, nil
}

func (usecase *bookmarkUsecase) ListBookmarksByUserID(ctx context.Context, UserID int) ([]domain.Bookmark, error) {
	bookmarks, err := usecase.repo.ListBookmarksByUserID(UserID)
	if err != nil {
		return []domain.Bookmark{}, status.Errorf(codes.Internal, "failed to get bookmark list: %v", err)
	}
	return *bookmarks, nil
}

func (usecase *bookmarkUsecase) ListBookmarksByArticleID(ctx context.Context, articleID int) ([]domain.Bookmark, error) {
	bookmarks, err := usecase.repo.ListBookmarksByArticleID(articleID)
	if err != nil {
		return []domain.Bookmark{}, status.Errorf(codes.Internal, "failed to get bookmark list: %v", err)
	}
	return *bookmarks, nil
}

func (usecase *bookmarkUsecase) DeleteBookmark(ctx context.Context, bookmark domain.Bookmark) (bool, error) {
	if err := usecase.repo.DeleteBookmarkByUserIDAndArticleID(int(bookmark.UserID), int(bookmark.ArticleID)); err != nil {
		return false, status.Errorf(codes.Internal, "failed to delete bookmark: %v", err)
	}
	return true, nil
}

func (usecase *bookmarkUsecase) DeleteBookmarkByUserID(ctx context.Context, userID int) (bool, error) {
	if err := usecase.repo.DeleteBookmarkByUserID(userID); err != nil {
		return false, status.Errorf(codes.Internal, "failed to delete bookmark: %v", err)
	}
	return true, nil
}

func (usecase *bookmarkUsecase) DeleteBookmarkByUserIDCompensate(ctx context.Context, userID int) (bool, error) {
	if err := usecase.repo.UpdateBookmarkByUserIDWithUnscoped(userID); err != nil {
		return false, status.Errorf(codes.Internal, "failed to compensate bookmark deletion: %v", err)
	}
	return true, nil
}

func (usecase *bookmarkUsecase) DeleteBookmarkByArticleID(ctx context.Context, articleID int) (bool, error) {
	if err := usecase.repo.DeleteBookmarkByArticleID(articleID); err != nil {
		return false, status.Errorf(codes.Internal, "failed to delete bookmark: %v", err)
	}
	return true, nil
}

func (usecase *bookmarkUsecase) DeleteBookmarkByArticleIDCompensate(ctx context.Context, articleID int) (bool, error) {
	if err := usecase.repo.UpdateBookmarkByArticleIDWithUnscoped(articleID); err != nil {
		return false, status.Errorf(codes.Internal, "failed to compensate bookmark deletion: %v", err)
	}
	return true, nil
}
