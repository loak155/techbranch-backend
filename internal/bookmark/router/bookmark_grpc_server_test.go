package router

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/loak155/techbranch-backend/internal/bookmark/domain"
	"github.com/loak155/techbranch-backend/internal/bookmark/usecase"
	"github.com/loak155/techbranch-backend/mock"
	articlepb "github.com/loak155/techbranch-backend/proto/article"
	bookmarkpb "github.com/loak155/techbranch-backend/proto/bookmark"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func TestCreateBookmark(t *testing.T) {
	type args struct {
		ctx context.Context
		req *bookmarkpb.CreateBookmarkRequest
	}

	req := &bookmarkpb.CreateBookmarkRequest{
		Bookmark: &bookmarkpb.Bookmark{
			UserId:    1,
			ArticleId: 1,
		},
	}

	resSuccessIncrementBookmarksCount := &articlepb.IncrementBookmarksCountResponse{
		Success: true,
	}

	resFalseIncrementBookmarksCount := &articlepb.IncrementBookmarksCountResponse{
		Success: false,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res *bookmarkpb.CreateBookmarkResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient) {
				mockArticleClient.EXPECT().IncrementBookmarksCount(gomock.Any(), gomock.Any()).Return(resSuccessIncrementBookmarksCount, nil)
				repo.EXPECT().GetBookmarkByUserIDAndArticleIDWithUnscoped(gomock.Any(), gomock.Any()).Return(&domain.Bookmark{}, nil)
				repo.EXPECT().CreateBookmark(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.CreateBookmarkResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, req.Bookmark.UserId, res.Bookmark.UserId)
				assert.Equal(t, req.Bookmark.ArticleId, res.Bookmark.ArticleId)
				assert.NotNil(t, res.Bookmark.CreatedAt)
				assert.NotNil(t, res.Bookmark.UpdatedAt)
			},
		},
		{
			name: "DuplicateEmail",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient) {
				mockArticleClient.EXPECT().IncrementBookmarksCount(gomock.Any(), gomock.Any()).Return(resFalseIncrementBookmarksCount, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.CreateBookmarkResponse, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := mock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.CreateBookmark(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestGetBookmark(t *testing.T) {
	type args struct {
		ctx context.Context
		req *bookmarkpb.GetBookmarkRequest
	}

	req := &bookmarkpb.GetBookmarkRequest{
		Id: 1,
	}

	repoResBookmark := &domain.Bookmark{
		ID:        1,
		UserID:    1,
		ArticleID: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, res *bookmarkpb.GetBookmarkResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().GetBookmark(gomock.Any()).Return(repoResBookmark, nil)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.GetBookmarkResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResBookmark.UserID, uint(res.Bookmark.UserId))
				assert.Equal(t, repoResBookmark.ArticleID, uint(res.Bookmark.ArticleId))
				assert.NotNil(t, res.Bookmark.CreatedAt)
				assert.NotNil(t, res.Bookmark.UpdatedAt)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().GetBookmark(gomock.Any()).Return(&domain.Bookmark{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.GetBookmarkResponse, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := mock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.GetBookmark(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestListBookmarks(t *testing.T) {
	type args struct {
		ctx context.Context
		req *bookmarkpb.ListBookmarksRequest
	}

	req := &bookmarkpb.ListBookmarksRequest{}

	repoResBookmarks := &[]domain.Bookmark{
		{
			ID:        1,
			UserID:    1,
			ArticleID: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
		{
			ID:        2,
			UserID:    2,
			ArticleID: 2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, res *bookmarkpb.ListBookmarksResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().ListBookmarks().Return(repoResBookmarks, nil)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.ListBookmarksResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(*repoResBookmarks), len(res.Bookmarks))
				for i, repoResBookmark := range *repoResBookmarks {
					assert.Equal(t, repoResBookmark.UserID, uint(res.Bookmarks[i].UserId))
					assert.Equal(t, repoResBookmark.ArticleID, uint(res.Bookmarks[i].ArticleId))
					assert.NotNil(t, res.Bookmarks[i].CreatedAt)
					assert.NotNil(t, res.Bookmarks[i].UpdatedAt)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := mock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.ListBookmarks(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestListBookmarksByUserID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *bookmarkpb.ListBookmarksByUserIDRequest
	}

	req := &bookmarkpb.ListBookmarksByUserIDRequest{
		UserId: 1,
	}

	repoResBookmarks := &[]domain.Bookmark{
		{
			ID:        1,
			UserID:    1,
			ArticleID: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
		{
			ID:        2,
			UserID:    2,
			ArticleID: 2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, res *bookmarkpb.ListBookmarksByUserIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().ListBookmarksByUserID(gomock.Any()).Return(repoResBookmarks, nil)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.ListBookmarksByUserIDResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(*repoResBookmarks), len(res.Bookmarks))
				for i, repoResBookmark := range *repoResBookmarks {
					assert.Equal(t, repoResBookmark.UserID, uint(res.Bookmarks[i].UserId))
					assert.Equal(t, repoResBookmark.ArticleID, uint(res.Bookmarks[i].ArticleId))
					assert.NotNil(t, res.Bookmarks[i].CreatedAt)
					assert.NotNil(t, res.Bookmarks[i].UpdatedAt)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := mock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
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
		req *bookmarkpb.ListBookmarksByArticleIDRequest
	}

	req := &bookmarkpb.ListBookmarksByArticleIDRequest{
		ArticleId: 1,
	}

	repoResBookmarks := &[]domain.Bookmark{
		{
			ID:        1,
			UserID:    1,
			ArticleID: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
		{
			ID:        2,
			UserID:    2,
			ArticleID: 2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, res *bookmarkpb.ListBookmarksByArticleIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().ListBookmarksByArticleID(gomock.Any()).Return(repoResBookmarks, nil)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.ListBookmarksByArticleIDResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(*repoResBookmarks), len(res.Bookmarks))
				for i, repoResBookmark := range *repoResBookmarks {
					assert.Equal(t, repoResBookmark.UserID, uint(res.Bookmarks[i].UserId))
					assert.Equal(t, repoResBookmark.ArticleID, uint(res.Bookmarks[i].ArticleId))
					assert.NotNil(t, res.Bookmarks[i].CreatedAt)
					assert.NotNil(t, res.Bookmarks[i].UpdatedAt)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := mock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.ListBookmarksByArticleID(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteBookmark(t *testing.T) {
	type args struct {
		ctx context.Context
		req *bookmarkpb.DeleteBookmarkRequest
	}

	req := &bookmarkpb.DeleteBookmarkRequest{
		Bookmark: &bookmarkpb.Bookmark{
			Id:        1,
			UserId:    1,
			ArticleId: 1,
		},
	}

	resSuccessDecrementBookmarksCount := &articlepb.DecrementBookmarksCountResponse{
		Success: true,
	}

	resFalseDecrementBookmarksCount := &articlepb.DecrementBookmarksCountResponse{
		Success: false,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res *bookmarkpb.DeleteBookmarkResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient) {
				mockArticleClient.EXPECT().DecrementBookmarksCount(gomock.Any(), gomock.Any()).Return(resSuccessDecrementBookmarksCount, nil)
				repo.EXPECT().DeleteBookmarkByUserIDAndArticleID(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.DeleteBookmarkResponse, err error) {
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
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient) {
				mockArticleClient.EXPECT().DecrementBookmarksCount(gomock.Any(), gomock.Any()).Return(resFalseDecrementBookmarksCount, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.DeleteBookmarkResponse, err error) {
				assert.Error(t, err)
				assert.False(t, res.Success)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := mock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.DeleteBookmark(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteBookmarkByUserID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *bookmarkpb.DeleteBookmarkByUserIDRequest
	}

	req := &bookmarkpb.DeleteBookmarkByUserIDRequest{
		UserId: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res *bookmarkpb.DeleteBookmarkByUserIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient) {
				repo.EXPECT().DeleteBookmarkByUserID(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.DeleteBookmarkByUserIDResponse, err error) {
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
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient) {
				repo.EXPECT().DeleteBookmarkByUserID(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.DeleteBookmarkByUserIDResponse, err error) {
				assert.Error(t, err)
				assert.False(t, res.Success)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := mock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.DeleteBookmarkByUserID(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteBookmarkByUserIDCompensate(t *testing.T) {
	type args struct {
		ctx context.Context
		req *bookmarkpb.DeleteBookmarkByUserIDRequest
	}

	req := &bookmarkpb.DeleteBookmarkByUserIDRequest{
		UserId: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res *bookmarkpb.DeleteBookmarkByUserIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient) {
				repo.EXPECT().UpdateBookmarkByUserIDWithUnscoped(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.DeleteBookmarkByUserIDResponse, err error) {
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
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient) {
				repo.EXPECT().UpdateBookmarkByUserIDWithUnscoped(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.DeleteBookmarkByUserIDResponse, err error) {
				assert.Error(t, err)
				assert.False(t, res.Success)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := mock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.DeleteBookmarkByUserIDCompensate(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteBookmarkByArticleID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *bookmarkpb.DeleteBookmarkByArticleIDRequest
	}

	req := &bookmarkpb.DeleteBookmarkByArticleIDRequest{
		ArticleId: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res *bookmarkpb.DeleteBookmarkByArticleIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient) {
				repo.EXPECT().DeleteBookmarkByArticleID(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.DeleteBookmarkByArticleIDResponse, err error) {
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
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient) {
				repo.EXPECT().DeleteBookmarkByArticleID(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.DeleteBookmarkByArticleIDResponse, err error) {
				assert.Error(t, err)
				assert.False(t, res.Success)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := mock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.DeleteBookmarkByArticleID(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteBookmarkByArticleIDCompensate(t *testing.T) {
	type args struct {
		ctx context.Context
		req *bookmarkpb.DeleteBookmarkByArticleIDRequest
	}

	req := &bookmarkpb.DeleteBookmarkByArticleIDRequest{
		ArticleId: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res *bookmarkpb.DeleteBookmarkByArticleIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient) {
				repo.EXPECT().UpdateBookmarkByArticleIDWithUnscoped(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.DeleteBookmarkByArticleIDResponse, err error) {
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
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *mock.MockArticleServiceClient) {
				repo.EXPECT().UpdateBookmarkByArticleIDWithUnscoped(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.DeleteBookmarkByArticleIDResponse, err error) {
				assert.Error(t, err)
				assert.False(t, res.Success)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := mock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase)
			res, err := s.DeleteBookmarkByArticleIDCompensate(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}
