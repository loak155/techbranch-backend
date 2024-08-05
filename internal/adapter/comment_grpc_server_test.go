package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/internal/usecase"
	"github.com/loak155/techbranch-backend/mock"
	"github.com/loak155/techbranch-backend/pkg/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func TestCreateComment(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CreateCommentRequest
	}

	req := &pb.CreateCommentRequest{
		UserId:    1,
		ArticleId: 1,
		Content:   "test_content",
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, res *pb.CreateCommentResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().CreateComment(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateCommentResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, req.UserId, res.Comment.UserId)
				assert.Equal(t, req.ArticleId, res.Comment.ArticleId)
				assert.NotNil(t, res.Comment.CreatedAt)
				assert.NotNil(t, res.Comment.UpdatedAt)
			},
		},
		{
			name: "InvalidArgument",
			args: args{
				ctx: context.Background(),
				req: &pb.CreateCommentRequest{},
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
			},
			checkResponse: func(t *testing.T, res *pb.CreateCommentResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid argument")
			},
		},
		{
			name: "InvalidData",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().CreateComment(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.CreateCommentResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to create comment")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewCommentUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewCommentGRPCServer(server, usecase)
			res, err := s.CreateComment(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestListCommentsByUserID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.ListCommentsByUserIDRequest
	}

	req := &pb.ListCommentsByUserIDRequest{
		UserId: 1,
	}

	repoResComments := []domain.Comment{
		{
			ID:        1,
			UserID:    1,
			ArticleID: 1,
			Content:   "test_content",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			UserID:    1,
			ArticleID: 2,
			Content:   "test_content2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, res *pb.ListCommentsByUserIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().ListCommentsByUserID(gomock.Any()).Return(&repoResComments, nil)
			},
			checkResponse: func(t *testing.T, res *pb.ListCommentsByUserIDResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(repoResComments), len(res.Comments))
				for i, repoResComment := range repoResComments {
					assert.Equal(t, repoResComment.UserID, uint(res.Comments[i].UserId))
					assert.Equal(t, repoResComment.ArticleID, uint(res.Comments[i].ArticleId))
					assert.NotNil(t, res.Comments[i].CreatedAt)
					assert.NotNil(t, res.Comments[i].UpdatedAt)
				}
			},
		},
		{
			name: "InvalidData",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().ListCommentsByUserID(gomock.Any()).Return(&[]domain.Comment{}, gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.ListCommentsByUserIDResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to list comments")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewCommentUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewCommentGRPCServer(server, usecase)
			res, err := s.ListCommentsByUserID(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestListCommentsByArticleID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.ListCommentsByArticleIDRequest
	}

	req := &pb.ListCommentsByArticleIDRequest{
		ArticleId: 1,
	}

	repoResComments := []domain.Comment{
		{
			ID:        1,
			UserID:    1,
			ArticleID: 1,
			Content:   "test_content",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			UserID:    2,
			ArticleID: 1,
			Content:   "test_content2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, res *pb.ListCommentsByArticleIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().ListCommentsByArticleID(gomock.Any()).Return(&repoResComments, nil)
			},
			checkResponse: func(t *testing.T, res *pb.ListCommentsByArticleIDResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(repoResComments), len(res.Comments))
				for i, repoResComment := range repoResComments {
					assert.Equal(t, repoResComment.UserID, uint(res.Comments[i].UserId))
					assert.Equal(t, repoResComment.ArticleID, uint(res.Comments[i].ArticleId))
					assert.NotNil(t, res.Comments[i].CreatedAt)
					assert.NotNil(t, res.Comments[i].UpdatedAt)
				}
			},
		},
		{
			name: "InvalidData",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().ListCommentsByArticleID(gomock.Any()).Return(&[]domain.Comment{}, gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.ListCommentsByArticleIDResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to list comments")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewCommentUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewCommentGRPCServer(server, usecase)
			res, err := s.ListCommentsByArticleID(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteComment(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.DeleteCommentRequest
	}

	req := &pb.DeleteCommentRequest{
		Id: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, res *pb.DeleteCommentResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteComment(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteCommentResponse, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteComment(gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteCommentResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete comment")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewCommentUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewCommentGRPCServer(server, usecase)
			res, err := s.DeleteComment(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteCommentByUserIDAndArticleID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.DeleteCommentByUserIDAndArticleIDRequest
	}

	req := &pb.DeleteCommentByUserIDAndArticleIDRequest{
		UserId:    1,
		ArticleId: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, res *pb.DeleteCommentByUserIDAndArticleIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByUserIDAndArticleID(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteCommentByUserIDAndArticleIDResponse, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByUserIDAndArticleID(gomock.Any(), gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteCommentByUserIDAndArticleIDResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete comment by user id and article id")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewCommentUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewCommentGRPCServer(server, usecase)
			res, err := s.DeleteCommentByUserIDAndArticleID(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteCommentByUserID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.DeleteCommentByUserIDRequest
	}

	req := &pb.DeleteCommentByUserIDRequest{
		UserId: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, res *pb.DeleteCommentByUserIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByUserID(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteCommentByUserIDResponse, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByUserID(gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteCommentByUserIDResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete comment by user id and article id")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewCommentUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewCommentGRPCServer(server, usecase)
			res, err := s.DeleteCommentByUserID(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteCommentByArticleID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.DeleteCommentByArticleIDRequest
	}

	req := &pb.DeleteCommentByArticleIDRequest{
		ArticleId: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, res *pb.DeleteCommentByArticleIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByArticleID(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteCommentByArticleIDResponse, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByArticleID(gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteCommentByArticleIDResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete comment by user id and article id")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewCommentUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewCommentGRPCServer(server, usecase)
			res, err := s.DeleteCommentByArticleID(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}
