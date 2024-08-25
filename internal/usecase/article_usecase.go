package usecase

import (
	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/internal/repository"
)

type IArticleUsecase interface {
	CreateArticle(article domain.Article) (domain.Article, error)
	GetArticle(id int) (domain.Article, error)
	ListArticles(offset, limit int) ([]domain.Article, error)
	UpdateArticle(article domain.Article) (domain.Article, error)
	DeleteArticle(id int) error
	GetArticleCount() (int, error)
}

type articleUsecase struct {
	repo repository.IArticleRepository
}

func NewArticleUsecase(repo repository.IArticleRepository) IArticleUsecase {
	return &articleUsecase{repo}
}

func (usecase *articleUsecase) CreateArticle(article domain.Article) (domain.Article, error) {
	if err := usecase.repo.CreateArticle(&article); err != nil {
		return domain.Article{}, err
	}
	return article, nil
}

func (usecase *articleUsecase) GetArticle(id int) (domain.Article, error) {
	article, err := usecase.repo.GetArticle(id)
	if err != nil {
		return domain.Article{}, err
	}
	return *article, nil
}

func (usecase *articleUsecase) ListArticles(offset, limit int) ([]domain.Article, error) {
	articles, err := usecase.repo.ListArticles(offset, limit)
	if err != nil {
		return []domain.Article{}, err
	}
	return *articles, nil
}

func (usecase *articleUsecase) UpdateArticle(article domain.Article) (domain.Article, error) {
	if err := usecase.repo.UpdateArticle(&article); err != nil {
		return domain.Article{}, err
	}
	return article, nil
}

func (usecase *articleUsecase) DeleteArticle(id int) error {
	err := usecase.repo.DeleteArticle(id)
	return err
}

func (usecase *articleUsecase) GetArticleCount() (int, error) {
	count, err := usecase.repo.GetArticleCount()
	if err != nil {
		return 0, err
	}
	return count, nil
}
