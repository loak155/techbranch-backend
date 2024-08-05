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

func TestCreateBookmark(t *testing.T) {
	type args struct {
		bookmark domain.Bookmark
	}

	reqBookmark := domain.Bookmark{
		UserID:    1,
		ArticleID: 1,
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
				bookmark: reqBookmark,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().CreateBookmark(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, resBookmark domain.Bookmark, err error) {
				assert.NoError(t, err)
				assert.Equal(t, reqBookmark.UserID, resBookmark.UserID)
				assert.Equal(t, reqBookmark.ArticleID, resBookmark.ArticleID)
				assert.NotNil(t, resBookmark.CreatedAt)
				assert.NotNil(t, resBookmark.UpdatedAt)
			},
		},
		{
			name: "InvalidData",
			args: args{
				bookmark: domain.Bookmark{},
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().CreateBookmark(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, resBookmark domain.Bookmark, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewBookmarkUsecase(repo)
			resUser, err := usecase.CreateBookmark(tc.args.bookmark)
			tc.checkResponse(t, resUser, err)
		})
	}
}

func TestGetBookmarkCountByArticleID(t *testing.T) {
	type args struct {
		articleID int
	}

	repoResBookmarkCount := 1

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, resBookmarkCount int, err error)
	}{
		{
			name: "OK",
			args: args{
				articleID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().GetBookmarkCountByArticleID(gomock.Any()).Return(repoResBookmarkCount, nil)
			},
			checkResponse: func(t *testing.T, resBookmarkCount int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResBookmarkCount, resBookmarkCount)
			},
		},
		{
			name: "NotFound",
			args: args{
				articleID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().GetBookmarkCountByArticleID(gomock.Any()).Return(0, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resBookmarkCount int, err error) {
				assert.Error(t, err)
				assert.Equal(t, 0, resBookmarkCount)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewBookmarkUsecase(repo)
			resBookmarkCount, err := usecase.GetBookmarkCountByArticleID(tc.args.articleID)
			tc.checkResponse(t, resBookmarkCount, err)
		})
	}
}

func TestListBookmarksByUserID(t *testing.T) {
	type args struct {
		userID int
	}

	repoResBookmarks := &[]domain.Bookmark{
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
		checkResponse func(t *testing.T, resBookmark []domain.Bookmark, err error)
	}{
		{
			name: "OK",
			args: args{
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
				}
			},
		},
		{
			name: "NotFound",
			args: args{},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().ListBookmarksByUserID(gomock.Any()).Return(&[]domain.Bookmark{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resBookmarks []domain.Bookmark, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewBookmarkUsecase(repo)
			resBookmarks, err := usecase.ListBookmarksByUserID(tc.args.userID)
			tc.checkResponse(t, resBookmarks, err)
		})
	}
}

func TestListBookmarksByArticleID(t *testing.T) {
	type args struct {
		articleID int
	}

	repoResBookmarks := &[]domain.Bookmark{
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
		checkResponse func(t *testing.T, resBookmark []domain.Bookmark, err error)
	}{
		{
			name: "OK",
			args: args{
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
				}
			},
		},
		{
			name: "NotFound",
			args: args{},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().ListBookmarksByArticleID(gomock.Any()).Return(&[]domain.Bookmark{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resBookmarks []domain.Bookmark, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewBookmarkUsecase(repo)
			resBookmarks, err := usecase.ListBookmarksByArticleID(tc.args.articleID)
			tc.checkResponse(t, resBookmarks, err)
		})
	}
}

func TestDeleteBookmarkByUserIDAndArticleID(t *testing.T) {
	type args struct {
		userID    int
		articleID int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			args: args{
				userID:    1,
				articleID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().DeleteBookmarkByUserIDAndArticleID(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			args: args{
				userID:    1,
				articleID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().DeleteBookmarkByUserIDAndArticleID(gomock.Any(), gomock.Any()).Return(gorm.ErrRecordNotFound)
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

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewBookmarkUsecase(repo)
			err := usecase.DeleteBookmarkByUserIDAndArticleID(tc.args.userID, tc.args.articleID)
			tc.checkResponse(t, err)
		})
	}
}

func TestDeleteBookmarkByUserID(t *testing.T) {
	type args struct {
		userID int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			args: args{
				userID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().DeleteBookmarkByUserID(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			args: args{
				userID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().DeleteBookmarkByUserID(gomock.Any()).Return(gorm.ErrRecordNotFound)
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

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewBookmarkUsecase(repo)
			err := usecase.DeleteBookmarkByUserID(tc.args.userID)
			tc.checkResponse(t, err)
		})
	}
}

func TestDeleteBookmarkByArticleID(t *testing.T) {
	type args struct {
		articleID int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIBookmarkRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			args: args{
				articleID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().DeleteBookmarkByArticleID(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			args: args{
				articleID: 1,
			},
			buildStubs: func(repo *mock.MockIBookmarkRepository) {
				repo.EXPECT().DeleteBookmarkByArticleID(gomock.Any()).Return(gorm.ErrRecordNotFound)
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

			repo := mock.NewMockIBookmarkRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewBookmarkUsecase(repo)
			err := usecase.DeleteBookmarkByArticleID(tc.args.articleID)
			tc.checkResponse(t, err)
		})
	}
}
