package router

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/loak155/techbranch-backend/pkg/config"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	articlemock "github.com/loak155/techbranch-backend/services/article/mock"
	"github.com/loak155/techbranch-backend/services/bookmark/internal/domain"
	"github.com/loak155/techbranch-backend/services/bookmark/internal/usecase"
	"github.com/loak155/techbranch-backend/services/bookmark/mock"
	bookmarkpb "github.com/loak155/techbranch-backend/services/bookmark/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func TestCreateBookmark(t *testing.T) {
	type args struct {
		ctx context.Context
		req *bookmarkpb.CreateBookmarkRequest
	}

	ctx := context.Background()

	req := &bookmarkpb.CreateBookmarkRequest{
		UserId:    1,
		ArticleId: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res *bookmarkpb.CreateBookmarkResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: ctx,
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().GetBookmarkByUserIDAndArticleIDWithUnscoped(gomock.Any(), gomock.Any()).Return(&domain.Bookmark{}, nil)
				repo.EXPECT().CreateBookmark(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.CreateBookmarkResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, req.UserId, res.Bookmark.UserId)
				assert.Equal(t, req.ArticleId, res.Bookmark.ArticleId)
				assert.NotNil(t, res.Bookmark.CreatedAt)
				assert.NotNil(t, res.Bookmark.UpdatedAt)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf, err := config.Load("../../../configs/config.yaml")
			assert.NoError(t, err)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			jwtManager := jwt.NewJwtManager(conf)
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase, jwtManager)
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
			conf, err := config.Load("../../../configs/config.yaml")
			assert.NoError(t, err)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			jwtManager := jwt.NewJwtManager(conf)
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase, jwtManager)
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
			conf, err := config.Load("../../../configs/config.yaml")
			assert.NoError(t, err)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			jwtManager := jwt.NewJwtManager(conf)
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase, jwtManager)
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
			conf, err := config.Load("../../../configs/config.yaml")
			assert.NoError(t, err)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			jwtManager := jwt.NewJwtManager(conf)
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase, jwtManager)
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
			conf, err := config.Load("../../../configs/config.yaml")
			assert.NoError(t, err)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			jwtManager := jwt.NewJwtManager(conf)
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase, jwtManager)
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

	ctx := context.Background()

	req := &bookmarkpb.DeleteBookmarkRequest{
		ArticleId: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res *bookmarkpb.DeleteBookmarkResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: ctx,
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
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
				ctx: ctx,
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().DeleteBookmarkByUserIDAndArticleID(gomock.Any(), gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *bookmarkpb.DeleteBookmarkResponse, err error) {
				assert.Error(t, err)
				assert.False(t, res.Success)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf, err := config.Load("../../../configs/config.yaml")
			assert.NoError(t, err)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			jwtManager := jwt.NewJwtManager(conf)
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase, jwtManager)
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
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res *bookmarkpb.DeleteBookmarkByUserIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
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
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
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
			conf, err := config.Load("../../../configs/config.yaml")
			assert.NoError(t, err)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			jwtManager := jwt.NewJwtManager(conf)
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase, jwtManager)
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
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res *bookmarkpb.DeleteBookmarkByUserIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
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
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
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
			conf, err := config.Load("../../../configs/config.yaml")
			assert.NoError(t, err)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			jwtManager := jwt.NewJwtManager(conf)
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase, jwtManager)
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
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res *bookmarkpb.DeleteBookmarkByArticleIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
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
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
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
			conf, err := config.Load("../../../configs/config.yaml")
			assert.NoError(t, err)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			jwtManager := jwt.NewJwtManager(conf)
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase, jwtManager)
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
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res *bookmarkpb.DeleteBookmarkByArticleIDResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
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
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
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
			conf, err := config.Load("../../../configs/config.yaml")
			assert.NoError(t, err)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := usecase.NewBookmarkUsecase(repo, mockArticleClient)
			server := grpc.NewServer()
			jwtManager := jwt.NewJwtManager(conf)
			server.GracefulStop()

			s := NewBookmarkGRPCServer(server, usecase, jwtManager)
			res, err := s.DeleteBookmarkByArticleIDCompensate(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}
