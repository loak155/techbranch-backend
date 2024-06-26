package repository

import (
	"github.com/loak155/techbranch-backend/services/article/internal/domain"
	"gorm.io/gorm"
)

type IArticleRepository interface {
	CreateArticle(article *domain.Article) error
	GetArticle(id int) (*domain.Article, error)
	ListArticles(offset, limit int) (*[]domain.Article, error)
	UpdateArticle(article *domain.Article) error
	DeleteArticle(id int) error
	GetArticleCount() (int, error)
	GetArticleByIDs(ids []int) (*[]domain.Article, error)
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) IArticleRepository {
	return &articleRepository{db}
}

func (repo *articleRepository) CreateArticle(article *domain.Article) error {
	err := repo.db.Create(article).Error
	return err
}

func (repo *articleRepository) GetArticle(id int) (*domain.Article, error) {
	article := &domain.Article{}
	err := repo.db.First(article, id).Error
	return article, err
}

func (repo *articleRepository) ListArticles(offset, limit int) (*[]domain.Article, error) {
	articles := &[]domain.Article{}
	err := repo.db.Order("created_at desc").Offset(offset).Limit(limit).Find(articles).Error
	return articles, err
}

func (repo *articleRepository) UpdateArticle(article *domain.Article) error {
	err := repo.db.Save(article).Error
	return err
}

func (repo *articleRepository) DeleteArticle(id int) error {
	err := repo.db.Delete(&domain.Article{}, id).Error
	return err
}

func (repo *articleRepository) GetArticleCount() (int, error) {
	var count int64
	err := repo.db.Model(&domain.Article{}).Count(&count).Error
	return int(count), err
}

func (repo *articleRepository) GetArticleByIDs(ids []int) (*[]domain.Article, error) {
	articles := &[]domain.Article{}
	err := repo.db.Where("id IN (?)", ids).Find(articles).Error
	return articles, err
}
