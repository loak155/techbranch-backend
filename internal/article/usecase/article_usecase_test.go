package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/loak155/techbranch-backend/internal/article/domain"
	"github.com/loak155/techbranch-backend/mock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateArticle(t *testing.T) {
	type args struct {
		ctx     context.Context
		article domain.Article
	}

	reqArticle := domain.Article{
		Title: "test_title",
		Url:   "https://example.com",
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, resArticle domain.Article, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:     context.Background(),
				article: reqArticle,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().CreateArticle(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, resArticle domain.Article, err error) {
				assert.NoError(t, err)
				assert.Equal(t, reqArticle.Title, resArticle.Title)
				assert.Equal(t, reqArticle.Url, resArticle.Url)
				assert.Equal(t, 0, int(resArticle.BookmarkCount))
				assert.NotNil(t, resArticle.CreatedAt)
				assert.NotNil(t, resArticle.UpdatedAt)
				assert.Equal(t, resArticle.DeletedAt, gorm.DeletedAt{})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIArticleRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewArticleUsecase(repo)
			resUser, err := usecase.CreateArticle(tc.args.ctx, tc.args.article)
			tc.checkResponse(t, resUser, err)
		})
	}
}

func TestGetArticle(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	repoResArticle := &domain.Article{
		ID:            1,
		Title:         "test_title",
		Url:           "https://example.com",
		BookmarkCount: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeletedAt:     gorm.DeletedAt{},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, resArticle domain.Article, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().GetArticle(gomock.Any()).Return(repoResArticle, nil)
			},
			checkResponse: func(t *testing.T, resArticle domain.Article, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResArticle.Title, resArticle.Title)
				assert.Equal(t, repoResArticle.Url, resArticle.Url)
				assert.Equal(t, repoResArticle.BookmarkCount, resArticle.BookmarkCount)
				assert.Equal(t, repoResArticle.CreatedAt, resArticle.CreatedAt)
				assert.Equal(t, repoResArticle.UpdatedAt, resArticle.UpdatedAt)
				assert.Equal(t, resArticle.DeletedAt, gorm.DeletedAt{})
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().GetArticle(gomock.Any()).Return(&domain.Article{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resArticle domain.Article, err error) {
				assert.Error(t, err)
				assert.Equal(t, resArticle, domain.Article{})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIArticleRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewArticleUsecase(repo)
			resArticle, err := usecase.GetArticle(tc.args.ctx, tc.args.id)
			tc.checkResponse(t, resArticle, err)
		})
	}
}

func TestListArticles(t *testing.T) {
	type args struct {
		ctx    context.Context
		offset int
		limit  int
	}

	repoResArticles := &[]domain.Article{
		{
			ID:            1,
			Title:         "test_Article1",
			Url:           "https://example1.com",
			BookmarkCount: 0,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     gorm.DeletedAt{},
		},
		{
			ID:            2,
			Title:         "test_Article2",
			Url:           "https://example2.com",
			BookmarkCount: 0,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     gorm.DeletedAt{},
		},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, resArticle []domain.Article, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:    context.Background(),
				offset: 0,
				limit:  10,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().ListArticles(gomock.Any(), gomock.Any()).Return(repoResArticles, nil)
			},
			checkResponse: func(t *testing.T, resArticles []domain.Article, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(*repoResArticles), len(resArticles))
				for i, repoResArticle := range *repoResArticles {
					assert.Equal(t, repoResArticle.Title, resArticles[i].Title)
					assert.Equal(t, repoResArticle.Url, resArticles[i].Url)
					assert.Equal(t, repoResArticle.BookmarkCount, resArticles[i].BookmarkCount)
					assert.Equal(t, repoResArticle.CreatedAt, resArticles[i].CreatedAt)
					assert.Equal(t, repoResArticle.UpdatedAt, resArticles[i].UpdatedAt)
					assert.Equal(t, repoResArticle.DeletedAt, resArticles[i].DeletedAt)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIArticleRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewArticleUsecase(repo)
			resArticles, err := usecase.ListArticles(tc.args.ctx, tc.args.offset, tc.args.limit)
			tc.checkResponse(t, resArticles, err)
		})
	}
}

func TestUpdateArticle(t *testing.T) {
	type args struct {
		ctx     context.Context
		article domain.Article
	}

	reqArticle := domain.Article{
		ID:            1,
		Title:         "test_title",
		Url:           "https://example.com",
		BookmarkCount: 0,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:     context.Background(),
				article: reqArticle,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().UpdateArticle(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "InvalidData",
			args: args{
				ctx:     context.Background(),
				article: reqArticle,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().UpdateArticle(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.Error(t, err)
				assert.False(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIArticleRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewArticleUsecase(repo)
			res, err := usecase.UpdateArticle(tc.args.ctx, tc.args.article)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteArticle(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().DeleteArticle(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().DeleteArticle(gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.Error(t, err)
				assert.False(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIArticleRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewArticleUsecase(repo)
			resArticle, err := usecase.DeleteArticle(tc.args.ctx, tc.args.id)
			tc.checkResponse(t, resArticle, err)
		})
	}
}

func TestIncrementBookmarksCount(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	repoResArticle := &domain.Article{
		ID:            1,
		Title:         "test_title",
		Url:           "https://example.com",
		BookmarkCount: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeletedAt:     gorm.DeletedAt{},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().GetArticle(gomock.Any()).Return(repoResArticle, nil)
				repo.EXPECT().UpdateArticle(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().GetArticle(gomock.Any()).Return(&domain.Article{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.Error(t, err)
				assert.False(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIArticleRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewArticleUsecase(repo)
			resArticle, err := usecase.IncrementBookmarksCount(tc.args.ctx, tc.args.id)
			tc.checkResponse(t, resArticle, err)
		})
	}
}

func TestIncrementBookmarksCountCompensate(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	repoResArticle := &domain.Article{
		ID:            1,
		Title:         "test_title",
		Url:           "https://example.com",
		BookmarkCount: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeletedAt:     gorm.DeletedAt{},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().GetArticle(gomock.Any()).Return(repoResArticle, nil)
				repo.EXPECT().UpdateArticle(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().GetArticle(gomock.Any()).Return(&domain.Article{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.Error(t, err)
				assert.False(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIArticleRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewArticleUsecase(repo)
			resArticle, err := usecase.IncrementBookmarksCountCompensate(tc.args.ctx, tc.args.id)
			tc.checkResponse(t, resArticle, err)
		})
	}
}

func TestDecrementBookmarksCount(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	repoResArticle := &domain.Article{
		ID:            1,
		Title:         "test_title",
		Url:           "https://example.com",
		BookmarkCount: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeletedAt:     gorm.DeletedAt{},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().GetArticle(gomock.Any()).Return(repoResArticle, nil)
				repo.EXPECT().UpdateArticle(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().GetArticle(gomock.Any()).Return(&domain.Article{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.Error(t, err)
				assert.False(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIArticleRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewArticleUsecase(repo)
			resArticle, err := usecase.DecrementBookmarksCount(tc.args.ctx, tc.args.id)
			tc.checkResponse(t, resArticle, err)
		})
	}
}

func TestDecrementBookmarksCountCompensate(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	repoResArticle := &domain.Article{
		ID:            1,
		Title:         "test_title",
		Url:           "https://example.com",
		BookmarkCount: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeletedAt:     gorm.DeletedAt{},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().GetArticle(gomock.Any()).Return(repoResArticle, nil)
				repo.EXPECT().UpdateArticle(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().GetArticle(gomock.Any()).Return(&domain.Article{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.Error(t, err)
				assert.False(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIArticleRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewArticleUsecase(repo)
			resArticle, err := usecase.DecrementBookmarksCountCompensate(tc.args.ctx, tc.args.id)
			tc.checkResponse(t, resArticle, err)
		})
	}
}
