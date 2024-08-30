package usecase

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/mock"
	"github.com/loak155/techbranch-backend/pkg/password"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	type args struct {
		user domain.User
	}

	reqUser := domain.User{
		Username: "test_username",
		Email:    "test@example.com",
		Password: "test_password",
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
				user: reqUser,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, resUser domain.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, reqUser.Username, resUser.Username)
				assert.NoError(t, password.CheckPassword(reqUser.Password, resUser.Password))
				assert.NotNil(t, resUser.CreatedAt)
				assert.NotNil(t, resUser.UpdatedAt)
			},
		},
		{
			name: "InvalidData",
			args: args{
				user: domain.User{},
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, resUser domain.User, err error) {
				assert.Error(t, err)
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
			resUser, err := usecase.CreateUser(tc.args.user)
			tc.checkResponse(t, resUser, err)
		})
	}
}

func TestGetUser(t *testing.T) {
	type args struct {
		id int
	}

	repoResUser := domain.User{
		ID:        1,
		Username:  "test_username",
		Email:     "test@example.com",
		Password:  "test_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
				id: 1,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUser(gomock.Any()).Return(&repoResUser, nil)
			},
			checkResponse: func(t *testing.T, resUser domain.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResUser.Username, resUser.Username)
				assert.Equal(t, repoResUser.Email, resUser.Email)
				assert.Equal(t, repoResUser.Password, resUser.Password)
				assert.Equal(t, repoResUser.CreatedAt, resUser.CreatedAt)
				assert.Equal(t, repoResUser.UpdatedAt, resUser.UpdatedAt)
			},
		},
		{
			name: "NotFound",
			args: args{
				id: 1,
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
			resUser, err := usecase.GetUser(tc.args.id)
			tc.checkResponse(t, resUser, err)
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	type args struct {
		email string
	}

	repoResUser := domain.User{
		ID:        1,
		Username:  "test_username",
		Email:     "test@example.com",
		Password:  "test_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
				email: "test@example.com",
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&repoResUser, nil)
			},
			checkResponse: func(t *testing.T, resUser domain.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResUser.Username, resUser.Username)
				assert.Equal(t, repoResUser.Email, resUser.Email)
				assert.Equal(t, repoResUser.Password, resUser.Password)
				assert.Equal(t, repoResUser.CreatedAt, resUser.CreatedAt)
				assert.Equal(t, repoResUser.UpdatedAt, resUser.UpdatedAt)
			},
		},
		{
			name: "NotFound",
			args: args{
				email: "notfound@example.com",
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
			resUser, err := usecase.GetUserByEmail(tc.args.email)
			tc.checkResponse(t, resUser, err)
		})
	}
}

func TestListUsers(t *testing.T) {
	type args struct {
		offset int
		limit  int
	}

	repoResUsers := []domain.User{
		{
			ID:        1,
			Username:  "test_username1",
			Email:     "test1@example.com",
			Password:  "test_password",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Username:  "test_username2",
			Email:     "test2@example.com",
			Password:  "test_password",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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
				offset: 0,
				limit:  10,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Return(&repoResUsers, nil)
			},
			checkResponse: func(t *testing.T, resUsers []domain.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(repoResUsers), len(resUsers))
				for i, repoResUser := range repoResUsers {
					assert.Equal(t, repoResUser.Username, resUsers[i].Username)
					assert.Equal(t, repoResUser.Email, resUsers[i].Email)
					assert.Equal(t, repoResUser.Password, resUsers[i].Password)
					assert.Equal(t, repoResUser.CreatedAt, resUsers[i].CreatedAt)
					assert.Equal(t, repoResUser.UpdatedAt, resUsers[i].UpdatedAt)
				}
			},
		},
		{
			name: "NotFound",
			args: args{
				offset: 0,
				limit:  10,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Return(&[]domain.User{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, resUsers []domain.User, err error) {
				assert.Error(t, err)
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
			resUsers, err := usecase.ListUsers(tc.args.offset, tc.args.limit)
			tc.checkResponse(t, resUsers, err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type args struct {
		user domain.User
	}

	reqUser := domain.User{
		ID:       1,
		Username: "test_username",
		Email:    "test@example.com",
		Password: "test_password",
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
				user: reqUser,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().UpdateUser(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, resUser domain.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, reqUser.Username, resUser.Username)
				assert.Equal(t, reqUser.Email, resUser.Email)
				assert.NoError(t, password.CheckPassword(reqUser.Password, resUser.Password))
			},
		},
		{
			name: "InvalidData",
			args: args{
				user: reqUser,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().UpdateUser(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, resUser domain.User, err error) {
				assert.Error(t, err)
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
			res, err := usecase.UpdateUser(tc.args.user)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type args struct {
		id int
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			args: args{
				id: 1,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().DeleteUser(gomock.Any()).Return(nil)
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
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().DeleteUser(gomock.Any()).Return(gorm.ErrRecordNotFound)
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

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := NewUserUsecase(repo)
			err := usecase.DeleteUser(tc.args.id)
			tc.checkResponse(t, err)
		})
	}
}
