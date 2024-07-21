package usecase

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/mock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateArticle(t *testing.T) {
	type args struct {
		article domain.Article
	}

	reqArticle := domain.Article{
		Title: "test_title",
		Url:   "https://example.com",
		Image: "http://example.com/image",
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
				article: reqArticle,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().CreateArticle(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, resArticle domain.Article, err error) {
				assert.NoError(t, err)
				assert.Equal(t, reqArticle.Title, resArticle.Title)
				assert.Equal(t, reqArticle.Url, resArticle.Url)
				assert.Equal(t, reqArticle.Image, resArticle.Image)
				assert.NotNil(t, resArticle.CreatedAt)
				assert.NotNil(t, resArticle.UpdatedAt)
			},
		},
		{
			name: "InvalidData",
			args: args{
				article: domain.Article{},
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().CreateArticle(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, resArticle domain.Article, err error) {
				assert.Error(t, err)
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
			resUser, err := usecase.CreateArticle(tc.args.article)
			tc.checkResponse(t, resUser, err)
		})
	}
}

func TestGetArticle(t *testing.T) {
	type args struct {
		id int
	}

	repoResArticle := domain.Article{
		ID:        1,
		Title:     "test_title",
		Url:       "https://example.com",
		Image:     "http://example.com/image",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
				id: 1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().GetArticle(gomock.Any()).Return(&repoResArticle, nil)
			},
			checkResponse: func(t *testing.T, resArticle domain.Article, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResArticle.Title, resArticle.Title)
				assert.Equal(t, repoResArticle.Url, resArticle.Url)
				assert.Equal(t, repoResArticle.Image, resArticle.Image)
				assert.Equal(t, repoResArticle.CreatedAt, resArticle.CreatedAt)
				assert.Equal(t, repoResArticle.UpdatedAt, resArticle.UpdatedAt)
			},
		},
		{
			name: "NotFound",
			args: args{
				id: 1,
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
			resArticle, err := usecase.GetArticle(tc.args.id)
			tc.checkResponse(t, resArticle, err)
		})
	}
}

func TestListArticles(t *testing.T) {
	type args struct {
		offset int
		limit  int
	}

	repoResArticles := []domain.Article{
		{
			ID:        1,
			Title:     "test_Article1",
			Url:       "https://example1.com",
			Image:     "http://example.com/image",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Title:     "test_Article2",
			Url:       "https://example2.com",
			Image:     "http://example2.com/image",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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
				offset: 0,
				limit:  10,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().ListArticles(gomock.Any(), gomock.Any()).Return(&repoResArticles, nil)
			},
			checkResponse: func(t *testing.T, resArticles []domain.Article, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(repoResArticles), len(resArticles))
				for i, repoResArticle := range repoResArticles {
					assert.Equal(t, repoResArticle.Title, resArticles[i].Title)
					assert.Equal(t, repoResArticle.Url, resArticles[i].Url)
					assert.Equal(t, repoResArticle.Image, resArticles[i].Image)
					assert.Equal(t, repoResArticle.CreatedAt, resArticles[i].CreatedAt)
					assert.Equal(t, repoResArticle.UpdatedAt, resArticles[i].UpdatedAt)
				}
			},
		},
		{
			name: "NotFound",
			args: args{
				offset: 0,
				limit:  10,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().ListArticles(gomock.Any(), gomock.Any()).Return(&[]domain.Article{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resArticles []domain.Article, err error) {
				assert.Error(t, err)
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
			resArticles, err := usecase.ListArticles(tc.args.offset, tc.args.limit)
			tc.checkResponse(t, resArticles, err)
		})
	}
}

func TestUpdateArticle(t *testing.T) {
	type args struct {
		article domain.Article
	}

	reqArticle := domain.Article{
		ID:    1,
		Title: "test_title",
		Url:   "https://example.com",
		Image: "http://example.com/image",
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
				article: reqArticle,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().UpdateArticle(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, resArticle domain.Article, err error) {
				assert.NoError(t, err)
				assert.Equal(t, reqArticle.Title, resArticle.Title)
				assert.Equal(t, reqArticle.Url, resArticle.Url)
				assert.Equal(t, reqArticle.Image, resArticle.Image)
			},
		},
		{
			name: "InvalidData",
			args: args{
				article: reqArticle,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().UpdateArticle(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, resArticle domain.Article, err error) {
				assert.Error(t, err)
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
			res, err := usecase.UpdateArticle(tc.args.article)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteArticle(t *testing.T) {
	type args struct {
		id int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIArticleRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			args: args{
				id: 1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().DeleteArticle(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			args: args{
				id: 1,
			},
			buildStubs: func(repo *mock.MockIArticleRepository) {
				repo.EXPECT().DeleteArticle(gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
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
			err := usecase.DeleteArticle(tc.args.id)
			tc.checkResponse(t, err)
		})
	}
}
