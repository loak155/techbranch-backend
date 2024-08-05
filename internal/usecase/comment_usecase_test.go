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

func TestCreateComment(t *testing.T) {
	type args struct {
		comment domain.Comment
	}

	reqComment := domain.Comment{
		UserID:    1,
		ArticleID: 1,
		Content:   "test_content",
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, resComment domain.Comment, err error)
	}{
		{
			name: "OK",
			args: args{
				comment: reqComment,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().CreateComment(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, resComment domain.Comment, err error) {
				assert.NoError(t, err)
				assert.Equal(t, reqComment.UserID, resComment.UserID)
				assert.Equal(t, reqComment.ArticleID, resComment.ArticleID)
				assert.NotNil(t, resComment.CreatedAt)
				assert.NotNil(t, resComment.UpdatedAt)
			},
		},
		{
			name: "InvalidData",
			args: args{
				comment: domain.Comment{},
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().CreateComment(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, resComment domain.Comment, err error) {
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

			usecase := NewCommentUsecase(repo)
			resUser, err := usecase.CreateComment(tc.args.comment)
			tc.checkResponse(t, resUser, err)
		})
	}
}

func TestListCommentsByUserID(t *testing.T) {
	type args struct {
		userID int
	}

	repoResComments := &[]domain.Comment{
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
		checkResponse func(t *testing.T, resComment []domain.Comment, err error)
	}{
		{
			name: "OK",
			args: args{
				userID: 1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().ListCommentsByUserID(gomock.Any()).Return(repoResComments, nil)
			},
			checkResponse: func(t *testing.T, resComments []domain.Comment, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(*repoResComments), len(resComments))
				for i, repoResComment := range *repoResComments {
					assert.Equal(t, repoResComment.UserID, resComments[i].UserID)
					assert.Equal(t, repoResComment.ArticleID, resComments[i].ArticleID)
					assert.Equal(t, repoResComment.CreatedAt, resComments[i].CreatedAt)
					assert.Equal(t, repoResComment.UpdatedAt, resComments[i].UpdatedAt)
				}
			},
		},
		{
			name: "NotFound",
			args: args{},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().ListCommentsByUserID(gomock.Any()).Return(&[]domain.Comment{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resComments []domain.Comment, err error) {
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

			usecase := NewCommentUsecase(repo)
			resComments, err := usecase.ListCommentsByUserID(tc.args.userID)
			tc.checkResponse(t, resComments, err)
		})
	}
}

func TestListCommentsByArticleID(t *testing.T) {
	type args struct {
		articleID int
	}

	repoResComments := &[]domain.Comment{
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
		checkResponse func(t *testing.T, resComment []domain.Comment, err error)
	}{
		{
			name: "OK",
			args: args{
				articleID: 1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().ListCommentsByArticleID(gomock.Any()).Return(repoResComments, nil)
			},
			checkResponse: func(t *testing.T, resComments []domain.Comment, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(*repoResComments), len(resComments))
				for i, repoResComment := range *repoResComments {
					assert.Equal(t, repoResComment.UserID, resComments[i].UserID)
					assert.Equal(t, repoResComment.ArticleID, resComments[i].ArticleID)
					assert.Equal(t, repoResComment.CreatedAt, resComments[i].CreatedAt)
					assert.Equal(t, repoResComment.UpdatedAt, resComments[i].UpdatedAt)
				}
			},
		},
		{
			name: "NotFound",
			args: args{},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().ListCommentsByArticleID(gomock.Any()).Return(&[]domain.Comment{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resComments []domain.Comment, err error) {
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

			usecase := NewCommentUsecase(repo)
			resComments, err := usecase.ListCommentsByArticleID(tc.args.articleID)
			tc.checkResponse(t, resComments, err)
		})
	}
}

func TestDeleteComment(t *testing.T) {
	type args struct {
		id int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			args: args{
				id: 1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteComment(gomock.Any()).Return(nil)
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
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteComment(gomock.Any()).Return(gorm.ErrRecordNotFound)
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

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewCommentUsecase(repo)
			err := usecase.DeleteComment(tc.args.id)
			tc.checkResponse(t, err)
		})
	}
}

func TestDeleteCommentByUserIDAndArticleID(t *testing.T) {
	type args struct {
		userID    int
		articleID int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			args: args{
				userID:    1,
				articleID: 1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByUserIDAndArticleID(gomock.Any(), gomock.Any()).Return(nil)
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
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByUserIDAndArticleID(gomock.Any(), gomock.Any()).Return(gorm.ErrRecordNotFound)
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

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewCommentUsecase(repo)
			err := usecase.DeleteCommentByUserIDAndArticleID(tc.args.userID, tc.args.articleID)
			tc.checkResponse(t, err)
		})
	}
}

func TestDeleteCommentByUserID(t *testing.T) {
	type args struct {
		userID int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			args: args{
				userID: 1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByUserID(gomock.Any()).Return(nil)
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
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByUserID(gomock.Any()).Return(gorm.ErrRecordNotFound)
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

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewCommentUsecase(repo)
			err := usecase.DeleteCommentByUserID(tc.args.userID)
			tc.checkResponse(t, err)
		})
	}
}

func TestDeleteCommentByArticleID(t *testing.T) {
	type args struct {
		articleID int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			args: args{
				articleID: 1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByArticleID(gomock.Any()).Return(nil)
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
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByArticleID(gomock.Any()).Return(gorm.ErrRecordNotFound)
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

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewCommentUsecase(repo)
			err := usecase.DeleteCommentByArticleID(tc.args.articleID)
			tc.checkResponse(t, err)
		})
	}
}
