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

func TestCreateBookmark(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CreateBookmarkRequest
	}

	req := &pb.CreateBookmarkRequest{
		UserId:    1,
		ArticleId: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, res *pb.CreateBookmarkResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().CreateBookmark(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateBookmarkResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, req.UserId, res.Bookmark.UserId)
				assert.Equal(t, req.ArticleId, res.Bookmark.ArticleId)
				assert.NotNil(t, res.Bookmark.CreatedAt)
				assert.NotNil(t, res.Bookmark.UpdatedAt)
			},
		},
		{
			name: "InvalidData",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().CreateBookmark(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.CreateBookmarkResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to create bookmark")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.CreateBookmark(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestGetBookmarkCountByArticleID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.GetBookmarkCountByArticleIDRequest
	}

	req := &pb.GetBookmarkCountByArticleIDRequest{
		ArticleId: 1,
	}

	repoResBookmarkCount := 1

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, res *pb.GetBookmarkCountByArticleIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().GetBookmarkCountByArticleID(gomock.Any()).Return(repoResBookmarkCount, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetBookmarkCountByArticleIDResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResBookmarkCount, int(res.Count))
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().GetBookmarkCountByArticleID(gomock.Any()).Return(0, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.GetBookmarkCountByArticleIDResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get bookmark count by article id")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.GetBookmarkCountByArticleID(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestListBookmarksByUserID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.ListBookmarksByUserIDRequest
	}

	req := &pb.ListBookmarksByUserIDRequest{
		UserId: 1,
	}

	repoResBookmarks := []domain.Bookmark{
		{
			ID:        1,
			UserID:    1,
			ArticleID: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			UserID:    1,
			ArticleID: 2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, res *pb.ListBookmarksByUserIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().ListBookmarksByUserID(gomock.Any()).Return(&repoResBookmarks, nil)
			},
			checkResponse: func(t *testing.T, res *pb.ListBookmarksByUserIDResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(repoResBookmarks), len(res.Bookmarks))
				for i, repoResBookmark := range repoResBookmarks {
					assert.Equal(t, repoResBookmark.UserID, uint(res.Bookmarks[i].UserId))
					assert.Equal(t, repoResBookmark.ArticleID, uint(res.Bookmarks[i].ArticleId))
					assert.NotNil(t, res.Bookmarks[i].CreatedAt)
					assert.NotNil(t, res.Bookmarks[i].UpdatedAt)
				}
			},
		},
		{
			name: "InvalidData",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().ListBookmarksByUserID(gomock.Any()).Return(&[]domain.Bookmark{}, gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.ListBookmarksByUserIDResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to list bookmarks")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.ListBookmarksByUserID(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestListBookmarksByArticleID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.ListBookmarksByArticleIDRequest
	}

	req := &pb.ListBookmarksByArticleIDRequest{
		ArticleId: 1,
	}

	repoResBookmarks := []domain.Bookmark{
		{
			ID:        1,
			UserID:    1,
			ArticleID: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			UserID:    2,
			ArticleID: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, res *pb.ListBookmarksByArticleIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().ListBookmarksByArticleID(gomock.Any()).Return(&repoResBookmarks, nil)
			},
			checkResponse: func(t *testing.T, res *pb.ListBookmarksByArticleIDResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(repoResBookmarks), len(res.Bookmarks))
				for i, repoResBookmark := range repoResBookmarks {
					assert.Equal(t, repoResBookmark.UserID, uint(res.Bookmarks[i].UserId))
					assert.Equal(t, repoResBookmark.ArticleID, uint(res.Bookmarks[i].ArticleId))
					assert.NotNil(t, res.Bookmarks[i].CreatedAt)
					assert.NotNil(t, res.Bookmarks[i].UpdatedAt)
				}
			},
		},
		{
			name: "InvalidData",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().ListBookmarksByArticleID(gomock.Any()).Return(&[]domain.Bookmark{}, gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.ListBookmarksByArticleIDResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to list bookmarks")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.ListBookmarksByArticleID(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteBookmarkByUserIDAndArticleID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.DeleteBookmarkByUserIDAndArticleIDRequest
	}

	req := &pb.DeleteBookmarkByUserIDAndArticleIDRequest{
		UserId:    1,
		ArticleId: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, res *pb.DeleteBookmarkByUserIDAndArticleIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().DeleteBookmarkByUserIDAndArticleID(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteBookmarkByUserIDAndArticleIDResponse, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().DeleteBookmarkByUserIDAndArticleID(gomock.Any(), gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteBookmarkByUserIDAndArticleIDResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete bookmark by user id and article id")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.DeleteBookmarkByUserIDAndArticleID(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteBookmarkByUserID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.DeleteBookmarkByUserIDRequest
	}

	req := &pb.DeleteBookmarkByUserIDRequest{
		UserId: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, res *pb.DeleteBookmarkByUserIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().DeleteBookmarkByUserID(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteBookmarkByUserIDResponse, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().DeleteBookmarkByUserID(gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteBookmarkByUserIDResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete bookmark by user id and article id")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.DeleteBookmarkByUserID(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteBookmarkByArticleID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.DeleteBookmarkByArticleIDRequest
	}

	req := &pb.DeleteBookmarkByArticleIDRequest{
		ArticleId: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, res *pb.DeleteBookmarkByArticleIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().DeleteBookmarkByArticleID(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteBookmarkByArticleIDResponse, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().DeleteBookmarkByArticleID(gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteBookmarkByArticleIDResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete bookmark by user id and article id")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.DeleteBookmarkByArticleID(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}
