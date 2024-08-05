package usecase

import (
	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/internal/repository"
)

type ICommentUsecase interface {
	CreateComment(comment domain.Comment) (domain.Comment, error)
	ListCommentsByUserID(userID int) ([]domain.Comment, error)
	ListCommentsByArticleID(articleID int) ([]domain.Comment, error)
	DeleteComment(id int) error
	DeleteCommentByUserIDAndArticleID(userID, articleID int) error
	DeleteCommentByUserID(userID int) error
	DeleteCommentByArticleID(articleID int) error
}

type commentUsecase struct {
	repo repository.ICommentRepository
}

func NewCommentUsecase(repo repository.ICommentRepository) ICommentUsecase {
	return &commentUsecase{repo}
}

func (usecase *commentUsecase) CreateComment(comment domain.Comment) (domain.Comment, error) {
	if err := usecase.repo.CreateComment(&comment); err != nil {
		return domain.Comment{}, err
	}
	return comment, nil
}

func (usecase *commentUsecase) ListCommentsByUserID(userID int) ([]domain.Comment, error) {
	comments, err := usecase.repo.ListCommentsByUserID(userID)
	if err != nil {
		return []domain.Comment{}, err
	}
	return *comments, nil
}

func (usecase *commentUsecase) ListCommentsByArticleID(articleID int) ([]domain.Comment, error) {
	comments, err := usecase.repo.ListCommentsByArticleID(articleID)
	if err != nil {
		return []domain.Comment{}, err
	}
	return *comments, nil
}

func (usecase *commentUsecase) DeleteComment(id int) error {
	err := usecase.repo.DeleteComment(id)
	return err
}

func (usecase *commentUsecase) DeleteCommentByUserIDAndArticleID(userID, articleID int) error {
	err := usecase.repo.DeleteCommentByUserIDAndArticleID(userID, articleID)
	return err
}

func (usecase *commentUsecase) DeleteCommentByUserID(userID int) error {
	err := usecase.repo.DeleteCommentByUserID(userID)
	return err
}

func (usecase *commentUsecase) DeleteCommentByArticleID(articleID int) error {
	err := usecase.repo.DeleteCommentByArticleID(articleID)
	return err
}
