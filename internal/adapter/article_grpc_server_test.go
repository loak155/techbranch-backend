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

func TestCreateArticle(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CreateArticleRequest
	}

	req := &pb.CreateArticleRequest{
		Title: "test_Article",
		Url:   "http://example.com",
		Image: "http://example.com/image",
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, res *pb.CreateArticleResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().CreateArticle(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateArticleResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, req.Title, res.Article.Title)
				assert.Equal(t, req.Url, res.Article.Url)
				assert.NotNil(t, res.Article.CreatedAt)
				assert.NotNil(t, res.Article.UpdatedAt)
			},
		},
		{
			name: "InvalidArgument",
			args: args{
				ctx: context.Background(),
				req: &pb.CreateArticleRequest{},
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
			},
			checkResponse: func(t *testing.T, res *pb.CreateArticleResponse, err error) {
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
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().CreateArticle(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.CreateArticleResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to create article")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIArticleRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewArticleUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewArticleGRPCServer(server, usecase)
			res, err := s.CreateArticle(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestGetArticle(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.GetArticleRequest
	}

	req := &pb.GetArticleRequest{
		Id: 1,
	}

	repoResArticle := domain.Article{
		ID:        1,
		Title:     "test_Article",
		Url:       "http://example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, res *pb.GetArticleResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().GetArticle(gomock.Any()).Return(&repoResArticle, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetArticleResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResArticle.Title, res.Article.Title)
				assert.Equal(t, repoResArticle.Url, res.Article.Url)
				assert.NotNil(t, res.Article.CreatedAt)
				assert.NotNil(t, res.Article.UpdatedAt)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().GetArticle(gomock.Any()).Return(&domain.Article{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.GetArticleResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get article")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIArticleRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewArticleUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewArticleGRPCServer(server, usecase)
			res, err := s.GetArticle(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestListArticles(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.ListArticlesRequest
	}

	req := &pb.ListArticlesRequest{
		Offset: 0,
		Limit:  10,
	}

	repoResArticles := []domain.Article{
		{

			ID:        1,
			Title:     "test_Article1",
			Url:       "http://example1.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{

			ID:        2,
			Title:     "test_Article2",
			Url:       "http://example2.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, res *pb.ListArticlesResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().ListArticles(gomock.Any(), gomock.Any()).Return(&repoResArticles, nil)
			},
			checkResponse: func(t *testing.T, res *pb.ListArticlesResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(repoResArticles), len(res.Articles))
				for i, repoResArticle := range repoResArticles {
					assert.Equal(t, repoResArticle.Title, res.Articles[i].Title)
					assert.Equal(t, repoResArticle.Url, res.Articles[i].Url)
					assert.NotNil(t, res.Articles[i].CreatedAt)
					assert.NotNil(t, res.Articles[i].UpdatedAt)
				}
			},
		},
		{
			name: "InvalidData",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().ListArticles(gomock.Any(), gomock.Any()).Return(&[]domain.Article{}, gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.ListArticlesResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to list articles")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIArticleRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewArticleUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewArticleGRPCServer(server, usecase)
			res, err := s.ListArticles(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestUpdateArticle(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateArticleRequest
	}

	req := &pb.UpdateArticleRequest{
		Id:    1,
		Title: "test_Article",
		Url:   "http://example.com",
		Image: "http://example.com/image",
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, res *pb.UpdateArticleResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().UpdateArticle(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateArticleResponse, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "InvalidArgument",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateArticleRequest{},
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
			},
			checkResponse: func(t *testing.T, res *pb.UpdateArticleResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid argument")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIArticleRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewArticleUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewArticleGRPCServer(server, usecase)
			res, err := s.UpdateArticle(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteArticle(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.DeleteArticleRequest
	}

	req := &pb.DeleteArticleRequest{
		Id: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, res *pb.DeleteArticleResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().DeleteArticle(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteArticleResponse, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().DeleteArticle(gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteArticleResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete article")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIArticleRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewArticleUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewArticleGRPCServer(server, usecase)
			res, err := s.DeleteArticle(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}
