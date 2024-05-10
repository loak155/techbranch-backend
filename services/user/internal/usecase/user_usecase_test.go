package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/loak155/techbranch-backend/pkg/password"
	"github.com/loak155/techbranch-backend/services/user/internal/domain"
	"github.com/loak155/techbranch-backend/services/user/mock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		user domain.User
	}

	reqUser := domain.User{
		Username: "test_user",
		Email:    "test@example.com",
		Password: "password",
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, resUser domain.User, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:  context.Background(),
				user: reqUser,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, resUser domain.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, reqUser.Username, resUser.Username)
				assert.Equal(t, reqUser.Email, resUser.Email)
				assert.NoError(t, password.CheckPassword(reqUser.Password, resUser.Password))
				assert.NotNil(t, resUser.CreatedAt)
				assert.NotNil(t, resUser.UpdatedAt)
				assert.Equal(t, resUser.DeletedAt, gorm.DeletedAt{})
			},
		},
		{
			name: "DuplicateEmail",
			args: args{
				ctx:  context.Background(),
				user: reqUser,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any()).Return(gorm.ErrDuplicatedKey)
			},
			checkResponse: func(t *testing.T, resUser domain.User, err error) {
				assert.Error(t, err)
				assert.Equal(t, resUser, domain.User{})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewUserUsecase(repo)
			resUser, err := usecase.CreateUser(tc.args.ctx, tc.args.user)
			tc.checkResponse(t, resUser, err)
		})
	}
}

func TestGetUser(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	repoResUser := &domain.User{
		ID:        1,
		Username:  "test_user",
		Email:     "test@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, resUser domain.User, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUser(gomock.Any()).Return(repoResUser, nil)
			},
			checkResponse: func(t *testing.T, resUser domain.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResUser.Username, resUser.Username)
				assert.Equal(t, repoResUser.Email, resUser.Email)
				assert.Equal(t, repoResUser.Password, resUser.Password)
				assert.Equal(t, repoResUser.CreatedAt, resUser.CreatedAt)
				assert.Equal(t, repoResUser.UpdatedAt, resUser.UpdatedAt)
				assert.Equal(t, resUser.DeletedAt, gorm.DeletedAt{})
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUser(gomock.Any()).Return(&domain.User{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resUser domain.User, err error) {
				assert.Error(t, err)
				assert.Equal(t, resUser, domain.User{})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewUserUsecase(repo)
			resUser, err := usecase.GetUser(tc.args.ctx, tc.args.id)
			tc.checkResponse(t, resUser, err)
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	type args struct {
		ctx   context.Context
		email string
	}

	repoResUser := &domain.User{
		ID:        1,
		Username:  "test_user",
		Email:     "test@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, resUser domain.User, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:   context.Background(),
				email: "test@example.com",
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(repoResUser, nil)
			},
			checkResponse: func(t *testing.T, resUser domain.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResUser.Username, resUser.Username)
				assert.Equal(t, repoResUser.Email, resUser.Email)
				assert.Equal(t, repoResUser.Password, resUser.Password)
				assert.Equal(t, repoResUser.CreatedAt, resUser.CreatedAt)
				assert.Equal(t, repoResUser.UpdatedAt, resUser.UpdatedAt)
				assert.Equal(t, resUser.DeletedAt, gorm.DeletedAt{})
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx:   context.Background(),
				email: "test@example.com",
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&domain.User{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resUser domain.User, err error) {
				assert.Error(t, err)
				assert.Equal(t, resUser, domain.User{})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewUserUsecase(repo)
			resUser, err := usecase.GetUserByEmail(tc.args.ctx, tc.args.email)
			tc.checkResponse(t, resUser, err)
		})
	}
}

func TestListUsers(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	repoResUsers := &[]domain.User{
		{

			ID:        1,
			Username:  "test_user",
			Email:     "test@example.com",
			Password:  "password",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
		{

			ID:        2,
			Username:  "test_user2",
			Email:     "test2@example.com",
			Password:  "password",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, resUser []domain.User, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().ListUsers().Return(repoResUsers, nil)
			},
			checkResponse: func(t *testing.T, resUsers []domain.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(*repoResUsers), len(resUsers))
				for i, repoResUser := range *repoResUsers {
					assert.Equal(t, repoResUser.Username, resUsers[i].Username)
					assert.Equal(t, repoResUser.Email, resUsers[i].Email)
					assert.Equal(t, repoResUser.Password, resUsers[i].Password)
					assert.Equal(t, repoResUser.CreatedAt, resUsers[i].CreatedAt)
					assert.Equal(t, repoResUser.UpdatedAt, resUsers[i].UpdatedAt)
					assert.Equal(t, repoResUser.DeletedAt, resUsers[i].DeletedAt)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewUserUsecase(repo)
			resUsers, err := usecase.ListUsers(tc.args.ctx)
			tc.checkResponse(t, resUsers, err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		user domain.User
	}

	reqUser := domain.User{
		ID:       1,
		Username: "test_user",
		Email:    "test@example.com",
		Password: "password",
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:  context.Background(),
				user: reqUser,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().UpdateUser(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res bool, err error) {
				assert.NoError(t, err)
				assert.True(t, res)
			},
		},
		{
			name: "InvalidData",
			args: args{
				ctx:  context.Background(),
				user: reqUser,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().UpdateUser(gomock.Any()).Return(gorm.ErrInvalidData)
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

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewUserUsecase(repo)
			res, err := usecase.UpdateUser(tc.args.ctx, tc.args.user)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, res bool, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().DeleteUser(gomock.Any()).Return(nil)
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
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().DeleteUser(gomock.Any()).Return(gorm.ErrRecordNotFound)
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

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewUserUsecase(repo)
			resUser, err := usecase.DeleteUser(tc.args.ctx, tc.args.id)
			tc.checkResponse(t, resUser, err)
		})
	}
}
