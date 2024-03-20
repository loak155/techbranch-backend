package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/loak155/techbranch-backend/internal/comment/domain"
	"github.com/loak155/techbranch-backend/mock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateComment(t *testing.T) {
	type args struct {
		ctx     context.Context
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
				ctx:     context.Background(),
				comment: reqComment,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().CreateComment(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, resComment domain.Comment, err error) {
				assert.NoError(t, err)
				assert.Equal(t, reqComment.UserID, resComment.UserID)
				assert.Equal(t, reqComment.ArticleID, resComment.ArticleID)
				assert.Equal(t, reqComment.Content, resComment.Content)
				assert.NotNil(t, resComment.CreatedAt)
				assert.NotNil(t, resComment.UpdatedAt)
				assert.Equal(t, resComment.DeletedAt, gorm.DeletedAt{})
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
			resComment, err := usecase.CreateComment(tc.args.ctx, tc.args.comment)
			tc.checkResponse(t, resComment, err)
		})
	}
}

func TestGetComment(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
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
		checkResponse func(t *testing.T, resComment domain.Comment, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().GetComment(gomock.Any()).Return(repoResComment, nil)
			},
			checkResponse: func(t *testing.T, resComment domain.Comment, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResComment.UserID, resComment.UserID)
				assert.Equal(t, repoResComment.ArticleID, resComment.ArticleID)
				assert.Equal(t, repoResComment.Content, resComment.Content)
				assert.Equal(t, repoResComment.CreatedAt, resComment.CreatedAt)
				assert.Equal(t, repoResComment.UpdatedAt, resComment.UpdatedAt)
				assert.Equal(t, resComment.DeletedAt, gorm.DeletedAt{})
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().GetComment(gomock.Any()).Return(&domain.Comment{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resComment domain.Comment, err error) {
				assert.Error(t, err)
				assert.Equal(t, resComment, domain.Comment{})
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
			resComment, err := usecase.GetComment(tc.args.ctx, tc.args.id)
			tc.checkResponse(t, resComment, err)
		})
	}
}

func TestListCommentsByArticleID(t *testing.T) {
	type args struct {
		ctx       context.Context
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
			DeletedAt: gorm.DeletedAt{},
		},
		{
			ID:        2,
			UserID:    2,
			ArticleID: 1,
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
		checkResponse func(t *testing.T, resComment []domain.Comment, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:       context.Background(),
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
					assert.Equal(t, repoResComment.Content, resComments[i].Content)
					assert.Equal(t, repoResComment.CreatedAt, resComments[i].CreatedAt)
					assert.Equal(t, repoResComment.UpdatedAt, resComments[i].UpdatedAt)
					assert.Equal(t, repoResComment.DeletedAt, resComments[i].DeletedAt)
				}
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
			resComments, err := usecase.ListCommentsByArticleID(tc.args.ctx, tc.args.articleID)
			tc.checkResponse(t, resComments, err)
		})
	}
}

func TestUpdateComment(t *testing.T) {
	type args struct {
		ctx     context.Context
		comment domain.Comment
	}

	reqComment := domain.Comment{
		ID:        1,
		UserID:    1,
		ArticleID: 1,
		Content:   "updated_content",
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
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:     context.Background(),
				comment: reqComment,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().GetComment(gomock.Any()).Return(repoResComment, nil)
				repo.EXPECT().UpdateComment(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx:     context.Background(),
				comment: reqComment,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().GetComment(gomock.Any()).Return(&domain.Comment{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.Error(t, err)
				assert.False(t, res)
			},
		},
		{
			name: "InvalidData",
			args: args{
				ctx:     context.Background(),
				comment: reqComment,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().GetComment(gomock.Any()).Return(repoResComment, nil)
				repo.EXPECT().UpdateComment(gomock.Any()).Return(gorm.ErrInvalidData)
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

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewCommentUsecase(repo)
			res, err := usecase.UpdateComment(tc.args.ctx, tc.args.comment)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteComment(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteComment(gomock.Any()).Return(nil)
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
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteComment(gomock.Any()).Return(gorm.ErrRecordNotFound)
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

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewCommentUsecase(repo)
			resComment, err := usecase.DeleteComment(tc.args.ctx, tc.args.id)
			tc.checkResponse(t, resComment, err)
		})
	}
}

func TestDeleteCommentByUserID(t *testing.T) {
	type args struct {
		ctx     context.Context
		user_id int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:     context.Background(),
				user_id: 1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByUserID(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx:     context.Background(),
				user_id: 1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByUserID(gomock.Any()).Return(gorm.ErrRecordNotFound)
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

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewCommentUsecase(repo)
			resComment, err := usecase.DeleteCommentByUserID(tc.args.ctx, tc.args.user_id)
			tc.checkResponse(t, resComment, err)
		})
	}
}

func TestDeleteCommentByUserIDCompensate(t *testing.T) {
	type args struct {
		ctx     context.Context
		user_id int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:     context.Background(),
				user_id: 1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().UpdateByUserIDWithUnscoped(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx:     context.Background(),
				user_id: 1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().UpdateByUserIDWithUnscoped(gomock.Any()).Return(gorm.ErrRecordNotFound)
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

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewCommentUsecase(repo)
			resComment, err := usecase.DeleteCommentByUserIDCompensate(tc.args.ctx, tc.args.user_id)
			tc.checkResponse(t, resComment, err)
		})
	}
}

func TestDeleteCommentByArticleID(t *testing.T) {
	type args struct {
		ctx        context.Context
		article_id int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:        context.Background(),
				article_id: 1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByArticleID(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx:        context.Background(),
				article_id: 1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().DeleteCommentByArticleID(gomock.Any()).Return(gorm.ErrRecordNotFound)
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

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewCommentUsecase(repo)
			resComment, err := usecase.DeleteCommentByArticleID(tc.args.ctx, tc.args.article_id)
			tc.checkResponse(t, resComment, err)
		})
	}
}

func TestDeleteCommentByArticleIDCompensate(t *testing.T) {
	type args struct {
		ctx        context.Context
		article_id int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockICommentRepository)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:        context.Background(),
				article_id: 1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().UpdateByArticleIDWithUnscoped(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx:        context.Background(),
				article_id: 1,
			},
			buildStubs: func(repo *mock.MockICommentRepository) {
				repo.EXPECT().UpdateByArticleIDWithUnscoped(gomock.Any()).Return(gorm.ErrRecordNotFound)
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

			repo := mock.NewMockICommentRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewCommentUsecase(repo)
			resComment, err := usecase.DeleteCommentByArticleIDCompensate(tc.args.ctx, tc.args.article_id)
			tc.checkResponse(t, resComment, err)
		})
	}
}
