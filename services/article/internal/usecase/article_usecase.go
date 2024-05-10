package usecase

import (
	"context"

	"github.com/loak155/techbranch-backend/services/article/internal/domain"
	"github.com/loak155/techbranch-backend/services/article/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IArticleUsecase interface {
	CreateArticle(ctx context.Context, article domain.Article) (domain.Article, error)
	GetArticle(ctx context.Context, id int) (domain.Article, error)
	ListArticles(ctx context.Context, offset, limit int) ([]domain.Article, error)
	UpdateArticle(ctx context.Context, article domain.Article) (bool, error)
	DeleteArticle(ctx context.Context, id int) (bool, error)
}

type articleUsecase struct {
	repo repository.IArticleRepository
}

func NewArticleUsecase(repo repository.IArticleRepository) IArticleUsecase {
	return &articleUsecase{repo}
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
