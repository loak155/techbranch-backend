package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	articlemock "github.com/loak155/techbranch-backend/services/article/mock"
	"github.com/loak155/techbranch-backend/services/bookmark/internal/domain"
	"github.com/loak155/techbranch-backend/services/bookmark/mock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateBookmark(t *testing.T) {
	type args struct {
		ctx      context.Context
		bookmark domain.Bookmark
	}

	reqBookmark := domain.Bookmark{
		UserID:    1,
		ArticleID: 1,
	}

	repoResBookmark := &domain.Bookmark{
		ID:        1,
		UserID:    1,
		ArticleID: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient)
		checkResponse func(t *testing.T, resBookmark domain.Bookmark, err error)
	}{
		{
			name: "Bookmark already exists, update it",
			args: args{
				ctx:      context.Background(),
				bookmark: reqBookmark,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().GetBookmarkByUserIDAndArticleIDWithUnscoped(gomock.Any(), gomock.Any()).Return(repoResBookmark, nil)
				repo.EXPECT().UpdateBookmarkWithUnscoped(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, resBookmark domain.Bookmark, err error) {
				assert.NoError(t, err)
				assert.Equal(t, reqBookmark.UserID, resBookmark.UserID)
				assert.Equal(t, reqBookmark.ArticleID, resBookmark.ArticleID)
				assert.NotNil(t, resBookmark.CreatedAt)
				assert.NotNil(t, resBookmark.UpdatedAt)
				assert.Equal(t, resBookmark.DeletedAt, gorm.DeletedAt{})
			},
		},
		{
			name: "Bookmark does not exist, create it",
			args: args{
				ctx:      context.Background(),
				bookmark: reqBookmark,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().GetBookmarkByUserIDAndArticleIDWithUnscoped(gomock.Any(), gomock.Any()).Return(&domain.Bookmark{}, nil)
				repo.EXPECT().CreateBookmark(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, resBookmark domain.Bookmark, err error) {
				assert.NoError(t, err)
				assert.Equal(t, reqBookmark.UserID, resBookmark.UserID)
				assert.Equal(t, reqBookmark.ArticleID, resBookmark.ArticleID)
				assert.NotNil(t, resBookmark.CreatedAt)
				assert.NotNil(t, resBookmark.UpdatedAt)
				assert.Equal(t, resBookmark.DeletedAt, gorm.DeletedAt{})
			},
		},
		{
			name: "Failed to update bookmark with unscoped",
			args: args{
				ctx:      context.Background(),
				bookmark: reqBookmark,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().GetBookmarkByUserIDAndArticleIDWithUnscoped(gomock.Any(), gomock.Any()).Return(repoResBookmark, nil)
				repo.EXPECT().UpdateBookmarkWithUnscoped(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, resBookmark domain.Bookmark, err error) {
				assert.Error(t, err)
				assert.Equal(t, resBookmark, domain.Bookmark{})
			},
		},
		{
			name: "Failed to create bookmark",
			args: args{
				ctx:      context.Background(),
				bookmark: reqBookmark,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().GetBookmarkByUserIDAndArticleIDWithUnscoped(gomock.Any(), gomock.Any()).Return(&domain.Bookmark{}, nil)
				repo.EXPECT().CreateBookmark(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, resBookmark domain.Bookmark, err error) {
				assert.Error(t, err)
				assert.Equal(t, resBookmark, domain.Bookmark{})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := NewBookmarkUsecase(repo, mockArticleClient)
			resBookmark, err := usecase.CreateBookmark(tc.args.ctx, tc.args.bookmark)
			tc.checkResponse(t, resBookmark, err)
		})
	}
}

func TestGetBookmark(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
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
		checkResponse func(t *testing.T, resBookmark domain.Bookmark, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().GetBookmark(gomock.Any()).Return(repoResBookmark, nil)
			},
			checkResponse: func(t *testing.T, resBookmark domain.Bookmark, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResBookmark.UserID, resBookmark.UserID)
				assert.Equal(t, repoResBookmark.ArticleID, resBookmark.ArticleID)
				assert.Equal(t, repoResBookmark.CreatedAt, resBookmark.CreatedAt)
				assert.Equal(t, repoResBookmark.UpdatedAt, resBookmark.UpdatedAt)
				assert.Equal(t, resBookmark.DeletedAt, gorm.DeletedAt{})
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().GetBookmark(gomock.Any()).Return(&domain.Bookmark{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resBookmark domain.Bookmark, err error) {
				assert.Error(t, err)
				assert.Equal(t, resBookmark, domain.Bookmark{})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewBookmarkUsecase(repo, mockArticleClient)
			resBookmark, err := usecase.GetBookmark(tc.args.ctx, tc.args.id)
			tc.checkResponse(t, resBookmark, err)
		})
	}
}

func TestListBookmarks(t *testing.T) {
	type args struct {
		ctx context.Context
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
		checkResponse func(t *testing.T, resBookmark []domain.Bookmark, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().ListBookmarks().Return(repoResBookmarks, nil)
			},
			checkResponse: func(t *testing.T, resBookmarks []domain.Bookmark, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(*repoResBookmarks), len(resBookmarks))
				for i, repoResBookmark := range *repoResBookmarks {
					assert.Equal(t, repoResBookmark.UserID, resBookmarks[i].UserID)
					assert.Equal(t, repoResBookmark.ArticleID, resBookmarks[i].ArticleID)
					assert.Equal(t, repoResBookmark.CreatedAt, resBookmarks[i].CreatedAt)
					assert.Equal(t, repoResBookmark.UpdatedAt, resBookmarks[i].UpdatedAt)
					assert.Equal(t, repoResBookmark.DeletedAt, resBookmarks[i].DeletedAt)
				}
			},
		},
		{
			name: "NG",
			args: args{
				ctx: context.Background(),
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().ListBookmarks().Return(&[]domain.Bookmark{}, gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, resBookmarks []domain.Bookmark, err error) {
				assert.Error(t, err)
				assert.Equal(t, resBookmarks, []domain.Bookmark{})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewBookmarkUsecase(repo, mockArticleClient)
			resBookmarks, err := usecase.ListBookmarks(tc.args.ctx)
			tc.checkResponse(t, resBookmarks, err)
		})
	}
}

func TestListBookmarksByUserID(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int
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
		checkResponse func(t *testing.T, resBookmark []domain.Bookmark, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().ListBookmarksByUserID(gomock.Any()).Return(repoResBookmarks, nil)
			},
			checkResponse: func(t *testing.T, resBookmarks []domain.Bookmark, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(*repoResBookmarks), len(resBookmarks))
				for i, repoResBookmark := range *repoResBookmarks {
					assert.Equal(t, repoResBookmark.UserID, resBookmarks[i].UserID)
					assert.Equal(t, repoResBookmark.ArticleID, resBookmarks[i].ArticleID)
					assert.Equal(t, repoResBookmark.CreatedAt, resBookmarks[i].CreatedAt)
					assert.Equal(t, repoResBookmark.UpdatedAt, resBookmarks[i].UpdatedAt)
					assert.Equal(t, repoResBookmark.DeletedAt, resBookmarks[i].DeletedAt)
				}
			},
		},
		{
			name: "NG",
			args: args{
				ctx: context.Background(),
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().ListBookmarksByUserID(gomock.Any()).Return(&[]domain.Bookmark{}, gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, resBookmarks []domain.Bookmark, err error) {
				assert.Error(t, err)
				assert.Equal(t, resBookmarks, []domain.Bookmark{})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewBookmarkUsecase(repo, mockArticleClient)
			resBookmarks, err := usecase.ListBookmarksByUserID(tc.args.ctx, 1)
			tc.checkResponse(t, resBookmarks, err)
		})
	}
}

func TestListBookmarksByArticleID(t *testing.T) {
	type args struct {
		ctx       context.Context
		articleID int
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
		checkResponse func(t *testing.T, resBookmark []domain.Bookmark, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:       context.Background(),
				articleID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().ListBookmarksByArticleID(gomock.Any()).Return(repoResBookmarks, nil)
			},
			checkResponse: func(t *testing.T, resBookmarks []domain.Bookmark, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(*repoResBookmarks), len(resBookmarks))
				for i, repoResBookmark := range *repoResBookmarks {
					assert.Equal(t, repoResBookmark.UserID, resBookmarks[i].UserID)
					assert.Equal(t, repoResBookmark.ArticleID, resBookmarks[i].ArticleID)
					assert.Equal(t, repoResBookmark.CreatedAt, resBookmarks[i].CreatedAt)
					assert.Equal(t, repoResBookmark.UpdatedAt, resBookmarks[i].UpdatedAt)
					assert.Equal(t, repoResBookmark.DeletedAt, resBookmarks[i].DeletedAt)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewBookmarkUsecase(repo, mockArticleClient)
			resBookmarks, err := usecase.ListBookmarksByArticleID(tc.args.ctx, 1)
			tc.checkResponse(t, resBookmarks, err)
		})
	}
}

func TestDeleteBookmark(t *testing.T) {
	type args struct {
		ctx      context.Context
		bookmark domain.Bookmark
	}

	reqBookmark := domain.Bookmark{
		ID:        1,
		UserID:    1,
		ArticleID: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:      context.Background(),
				bookmark: reqBookmark,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().DeleteBookmarkByUserIDAndArticleID(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "Failed to delete bookmark",
			args: args{
				ctx:      context.Background(),
				bookmark: reqBookmark,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().DeleteBookmarkByUserIDAndArticleID(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidData)
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

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := NewBookmarkUsecase(repo, mockArticleClient)
			resBookmark, err := usecase.DeleteBookmark(tc.args.ctx, tc.args.bookmark)
			tc.checkResponse(t, resBookmark, err)
		})
	}
}

func TestDeleteBookmarkByUserID(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().DeleteBookmarkByUserID(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "NG",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().DeleteBookmarkByUserID(gomock.Any()).Return(gorm.ErrInvalidData)
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

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := NewBookmarkUsecase(repo, mockArticleClient)
			resBookmark, err := usecase.DeleteBookmarkByUserID(tc.args.ctx, tc.args.userID)
			tc.checkResponse(t, resBookmark, err)
		})
	}
}

func TestDeleteBookmarkByUserIDCompensate(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().UpdateBookmarkByUserIDWithUnscoped(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "NG",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().UpdateBookmarkByUserIDWithUnscoped(gomock.Any()).Return(gorm.ErrInvalidData)
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

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := NewBookmarkUsecase(repo, mockArticleClient)
			resBookmark, err := usecase.DeleteBookmarkByUserIDCompensate(tc.args.ctx, tc.args.userID)
			tc.checkResponse(t, resBookmark, err)
		})
	}
}

func TestDeleteBookmarkByArticleID(t *testing.T) {
	type args struct {
		ctx       context.Context
		articleID int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:       context.Background(),
				articleID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().DeleteBookmarkByArticleID(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "NG",
			args: args{
				ctx:       context.Background(),
				articleID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().DeleteBookmarkByArticleID(gomock.Any()).Return(gorm.ErrInvalidData)
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

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := NewBookmarkUsecase(repo, mockArticleClient)
			resBookmark, err := usecase.DeleteBookmarkByArticleID(tc.args.ctx, tc.args.articleID)
			tc.checkResponse(t, resBookmark, err)
		})
	}
}

func TestDeleteBookmarkByArticleIDCompensate(t *testing.T) {
	type args struct {
		ctx       context.Context
		articleID int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:       context.Background(),
				articleID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().UpdateBookmarkByArticleIDWithUnscoped(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "NG",
			args: args{
				ctx:       context.Background(),
				articleID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository, mockArticleClient *articlemock.MockArticleServiceClient) {
				repo.EXPECT().UpdateBookmarkByArticleIDWithUnscoped(gomock.Any()).Return(gorm.ErrInvalidData)
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

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			mockArticleClient := articlemock.NewMockArticleServiceClient(mockCtrl)
			tc.buildStubs(repo, mockArticleClient)

			usecase := NewBookmarkUsecase(repo, mockArticleClient)
			resBookmark, err := usecase.DeleteBookmarkByArticleIDCompensate(tc.args.ctx, tc.args.articleID)
			tc.checkResponse(t, resBookmark, err)
		})
	}
}
