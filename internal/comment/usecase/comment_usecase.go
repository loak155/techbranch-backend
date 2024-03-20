package usecase

import (
	"context"

	"github.com/loak155/techbranch-backend/internal/comment/domain"
	"github.com/loak155/techbranch-backend/internal/comment/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ICommentUsecase interface {
	CreateComment(ctx context.Context, comment domain.Comment) (domain.Comment, error)
	GetComment(ctx context.Context, id int) (domain.Comment, error)
	ListCommentsByArticleID(ctx context.Context, articleID int) ([]domain.Comment, error)
	UpdateComment(ctx context.Context, comment domain.Comment) (bool, error)
	DeleteComment(ctx context.Context, id int) (bool, error)
	DeleteCommentByUserID(ctx context.Context, userID int) (bool, error)
	DeleteCommentByUserIDCompensate(ctx context.Context, userID int) (bool, error)
	DeleteCommentByArticleID(ctx context.Context, articleID int) (bool, error)
	DeleteCommentByArticleIDCompensate(ctx context.Context, articleID int) (bool, error)
}

type commentUsecase struct {
	repo repository.ICommentRepository
}

func NewCommentUsecase(ur repository.ICommentRepository) ICommentUsecase {
	return &commentUsecase{ur}
}

func (usecase *commentUsecase) CreateComment(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
	if err := usecase.repo.CreateComment(&comment); err != nil {
		return domain.Comment{}, status.Errorf(codes.Internal, "failed to create comment: %v", err)
	}
	return comment, nil
}

func (usecase *commentUsecase) GetComment(ctx context.Context, id int) (domain.Comment, error) {
	comment, err := usecase.repo.GetComment(id)
	if err != nil {
		return domain.Comment{}, status.Errorf(codes.Internal, "failed to get comment: %v", err)
	}
	return *comment, nil
}

func (usecase *commentUsecase) ListCommentsByArticleID(ctx context.Context, articleID int) ([]domain.Comment, error) {
	comments, err := usecase.repo.ListCommentsByArticleID(articleID)
	if err != nil {
		return []domain.Comment{}, status.Errorf(codes.Internal, "failed to get comment list: %v", err)
	}
	return *comments, nil
}

func (usecase *commentUsecase) UpdateComment(ctx context.Context, comment domain.Comment) (bool, error) {
	storedComment, err := usecase.repo.GetComment(int(comment.ID))
	if err != nil {
		return false, status.Errorf(codes.Internal, "failed to get comment: %v", err)
	}

	storedComment.Content = comment.Content

	if err := usecase.repo.UpdateComment(storedComment); err != nil {
		return false, status.Errorf(codes.Internal, "failed to update comment: %v", err)
	}
	return true, nil
}

func (usecase *commentUsecase) DeleteComment(ctx context.Context, id int) (bool, error) {
	if err := usecase.repo.DeleteComment(id); err != nil {
		return false, status.Errorf(codes.Internal, "failed to delete comment: %v", err)
	}
	return true, nil
}

func (usecase *commentUsecase) DeleteCommentByUserID(ctx context.Context, userID int) (bool, error) {
	if err := usecase.repo.DeleteCommentByUserID(userID); err != nil {
		return false, status.Errorf(codes.Internal, "failed to delete comment: %v", err)
	}
	return true, nil
}

func (usecase *commentUsecase) DeleteCommentByUserIDCompensate(ctx context.Context, userID int) (bool, error) {
	if err := usecase.repo.UpdateByUserIDWithUnscoped(userID); err != nil {
		return false, status.Errorf(codes.Internal, "failed to compensate comment deletion: %v", err)
	}
	return true, nil
}

func (usecase *commentUsecase) DeleteCommentByArticleID(ctx context.Context, articleID int) (bool, error) {
	if err := usecase.repo.DeleteCommentByArticleID(articleID); err != nil {
		return false, status.Errorf(codes.Internal, "failed to delete comment: %v", err)
	}
	return true, nil
}

func (usecase *commentUsecase) DeleteCommentByArticleIDCompensate(ctx context.Context, articleID int) (bool, error) {
	if err := usecase.repo.UpdateByArticleIDWithUnscoped(articleID); err != nil {
		return false, status.Errorf(codes.Internal, "failed to compensate comment deletion: %v", err)
	}
	return true, nil
}
