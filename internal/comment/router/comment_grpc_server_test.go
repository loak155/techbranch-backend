package router

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/loak155/techbranch-backend/internal/comment/domain"
	"github.com/loak155/techbranch-backend/internal/comment/usecase"
	"github.com/loak155/techbranch-backend/mock"
	pb "github.com/loak155/techbranch-backend/proto/comment"
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
		Comment: &pb.Comment{
			UserId:    1,
			ArticleId: 1,
			Content:   "test_content",
		},
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
				assert.Equal(t, req.Comment.UserId, res.Comment.UserId)
				assert.Equal(t, req.Comment.ArticleId, res.Comment.ArticleId)
				assert.Equal(t, req.Comment.Content, res.Comment.Content)
				assert.NotNil(t, res.Comment.CreatedAt)
				assert.NotNil(t, res.Comment.UpdatedAt)
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
				assert.Nil(t, res)
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

func TestGetComment(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.GetCommentRequest
	}

	req := &pb.GetCommentRequest{
		Id: 1,
	}

	repoResComment := &domain.Comment{
		ID:        1,
		UserID:    1,
		ArticleID: 1,
		Content:   "test_content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, res *pb.GetCommentResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().GetComment(gomock.Any()).Return(repoResComment, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetCommentResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResComment.UserID, uint(res.Comment.UserId))
				assert.Equal(t, repoResComment.ArticleID, uint(res.Comment.ArticleId))
				assert.Equal(t, repoResComment.Content, res.Comment.Content)
				assert.NotNil(t, res.Comment.CreatedAt)
				assert.NotNil(t, res.Comment.UpdatedAt)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().GetComment(gomock.Any()).Return(&domain.Comment{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.GetCommentResponse, err error) {
				assert.Error(t, err)
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
			res, err := s.GetComment(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestListCommentsByArticleID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.ListCommentsByArticleIDRequest
	}

	req := &pb.ListCommentsByArticleIDRequest{}

	repoResComments := &[]domain.Comment{
		{
			ID:        1,
			UserID:    1,
			ArticleID: 1,
			Content:   "test_content1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
		{
			ID:        2,
			UserID:    1,
			ArticleID: 2,
			Content:   "test_content2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
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
				repo.EXPECT().ListCommentsByArticleID(gomock.Any()).Return(repoResComments, nil)
			},
			checkResponse: func(t *testing.T, res *pb.ListCommentsByArticleIDResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(*repoResComments), len(res.Comments))
				for i, repoResComment := range *repoResComments {
					assert.Equal(t, repoResComment.UserID, uint(res.Comments[i].UserId))
					assert.Equal(t, repoResComment.ArticleID, uint(res.Comments[i].ArticleId))
					assert.Equal(t, repoResComment.Content, res.Comments[i].Content)
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
				repo.EXPECT().ListCommentsByArticleID(gomock.Any()).Return(nil, gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.ListCommentsByArticleIDResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
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

func TestUpdateComment(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateCommentRequest
	}

	req := &pb.UpdateCommentRequest{
		Comment: &pb.Comment{
			Id:        1,
			UserId:    1,
			ArticleId: 1,
			Content:   "test_content",
		},
	}

	repoResComment := &domain.Comment{
		ID:        1,
		UserID:    1,
		ArticleID: 1,
		Content:   "test_content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, res *pb.UpdateCommentResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().GetComment(gomock.Any()).Return(repoResComment, nil)
				repo.EXPECT().UpdateComment(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateCommentResponse, err error) {
				assert.NoError(t, err)
				assert.True(t, res.Success)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().GetComment(gomock.Any()).Return(&domain.Comment{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateCommentResponse, err error) {
				assert.Error(t, err)
				assert.False(t, res.Success)
			},
		},
		{
			name: "InvalidData",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().GetComment(gomock.Any()).Return(repoResComment, nil)
				repo.EXPECT().UpdateComment(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateCommentResponse, err error) {
				assert.Error(t, err)
				assert.False(t, res.Success)
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
			res, err := s.UpdateComment(tc.args.ctx, tc.args.req)
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
				assert.True(t, res.Success)
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
				assert.False(t, res.Success)
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
				assert.True(t, res.Success)
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
				assert.False(t, res.Success)
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

func TestDeleteCommentByUserIDCompensate(t *testing.T) {
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
				repo.EXPECT().UpdateByUserIDWithUnscoped(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteCommentByUserIDResponse, err error) {
				assert.NoError(t, err)
				assert.True(t, res.Success)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().UpdateByUserIDWithUnscoped(gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteCommentByUserIDResponse, err error) {
				assert.Error(t, err)
				assert.False(t, res.Success)
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
			res, err := s.DeleteCommentByUserIDCompensate(tc.args.ctx, tc.args.req)
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
				assert.True(t, res.Success)
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
				assert.False(t, res.Success)
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

func TestDeleteCommentByArticleIDCompensate(t *testing.T) {
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
				repo.EXPECT().UpdateByArticleIDWithUnscoped(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteCommentByArticleIDResponse, err error) {
				assert.NoError(t, err)
				assert.True(t, res.Success)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().UpdateByArticleIDWithUnscoped(gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteCommentByArticleIDResponse, err error) {
				assert.Error(t, err)
				assert.False(t, res.Success)
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
			res, err := s.DeleteCommentByArticleIDCompensate(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}
