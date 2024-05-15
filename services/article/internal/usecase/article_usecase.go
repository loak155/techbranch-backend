package usecase

import (
	"context"

	"github.com/loak155/techbranch-backend/services/article/internal/domain"
	"github.com/loak155/techbranch-backend/services/article/internal/repository"
	pb "github.com/loak155/techbranch-backend/services/bookmark/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IArticleUsecase interface {
	CreateArticle(ctx context.Context, article domain.Article) (domain.Article, error)
	GetArticle(ctx context.Context, id int) (domain.Article, error)
	ListArticles(ctx context.Context, offset, limit int) ([]domain.Article, error)
	UpdateArticle(ctx context.Context, article domain.Article) (bool, error)
	DeleteArticle(ctx context.Context, id int) (bool, error)
	GetArticleCount(ctx context.Context) (int, error)
	GetBookmarkedArticle(ctx context.Context, userId int) ([]domain.Article, error)
}

type articleUsecase struct {
	repo           repository.IArticleRepository
	bookmarkClient pb.BookmarkServiceClient
}

func NewArticleUsecase(repo repository.IArticleRepository, bookmarkClient pb.BookmarkServiceClient) IArticleUsecase {
	return &articleUsecase{repo, bookmarkClient}
}

func (usecase *articleUsecase) CreateArticle(ctx context.Context, article domain.Article) (domain.Article, error) {
	if err := usecase.repo.CreateArticle(&article); err != nil {
		return domain.Article{}, status.Errorf(codes.Internal, "failed to create article: %v", err)
	}
	return article, nil
}

func (usecase *articleUsecase) GetArticle(ctx context.Context, id int) (domain.Article, error) {
	article, err := usecase.repo.GetArticle(id)
	if err != nil {
		return domain.Article{}, status.Errorf(codes.Internal, "failed to get article: %v", err)
	}
	return *article, nil
}

func (usecase *articleUsecase) ListArticles(ctx context.Context, offset, limit int) ([]domain.Article, error) {
	articles, err := usecase.repo.ListArticles(offset, limit)
	if err != nil {
		return []domain.Article{}, status.Errorf(codes.Internal, "failed to get article list: %v", err)
	}
	return *articles, nil
}

func (usecase *articleUsecase) UpdateArticle(ctx context.Context, article domain.Article) (bool, error) {
	updatedArticle := domain.Article{
		ID:    article.ID,
		Title: article.Title,
		Url:   article.Url,
		Image: article.Image,
	}
	if err := usecase.repo.UpdateArticle(&updatedArticle); err != nil {
		return false, status.Errorf(codes.Internal, "failed to update article: %v", err)
	}
	return true, nil
}

func (usecase *articleUsecase) DeleteArticle(ctx context.Context, id int) (bool, error) {
	if err := usecase.repo.DeleteArticle(id); err != nil {
		return false, status.Errorf(codes.Internal, "failed to delete article: %v", err)
	}
	return true, nil
}

func (usecase *articleUsecase) GetArticleCount(ctx context.Context) (int, error) {
	count, err := usecase.repo.GetArticleCount()
	if err != nil {
		return 0, status.Errorf(codes.Internal, "failed to get article count: %v", err)
	}
	return count, nil
}

func (usecase *articleUsecase) GetBookmarkedArticle(ctx context.Context, userId int) ([]domain.Article, error) {
	bookmarks, err := usecase.bookmarkClient.ListBookmarksByUserID(ctx, &pb.ListBookmarksByUserIDRequest{UserId: int32(userId)})
	if err != nil {
		return []domain.Article{}, status.Errorf(codes.Internal, "failed to get bookmarked article: %v", err)
	}

	articleIds := make([]int, 0)
	for _, bookmark := range bookmarks.Bookmarks {
		articleIds = append(articleIds, int(bookmark.ArticleId))
	}

	articles, err := usecase.repo.GetArticleByIDs(articleIds)
	if err != nil {
		return []domain.Article{}, status.Errorf(codes.Internal, "failed to get bookmarked article: %v", err)
	}
	return *articles, nil
}
