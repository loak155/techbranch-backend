package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/internal/usecase"
	"github.com/loak155/techbranch-backend/mock"
	"github.com/loak155/techbranch-backend/pkg/password"
	"github.com/loak155/techbranch-backend/pkg/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CreateUserRequest
	}

	req := &pb.CreateUserRequest{
		Username: "test_username",
		Email:    "test@example.com",
		Password: "test_password",
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, req.Username, res.User.Username)
				assert.Equal(t, req.Email, res.User.Email)
				assert.NoError(t, password.CheckPassword(req.Password, res.User.Password))
				assert.NotNil(t, res.User.CreatedAt)
				assert.NotNil(t, res.User.UpdatedAt)
			},
		},
		{
			name: "InvalidArgument",
			args: args{
				ctx: context.Background(),
				req: &pb.CreateUserRequest{},
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
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
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to create user")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewUserUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewUserGRPCServer(server, usecase)
			res, err := s.CreateUser(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestGetUser(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.GetUserRequest
	}

	req := &pb.GetUserRequest{
		Id: 1,
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
		checkResponse func(t *testing.T, res *pb.GetUserResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUser(gomock.Any()).Return(&repoResUser, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetUserResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResUser.Username, res.User.Username)
				assert.Equal(t, repoResUser.Email, res.User.Email)
				assert.Equal(t, repoResUser.Password, res.User.Password)
				assert.NotNil(t, res.User.CreatedAt)
				assert.NotNil(t, res.User.UpdatedAt)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUser(gomock.Any()).Return(&domain.User{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.GetUserResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get user")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewUserUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewUserGRPCServer(server, usecase)
			res, err := s.GetUser(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.GetUserByEmailRequest
	}

	req := &pb.GetUserByEmailRequest{
		Email: "test@example.com",
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
		checkResponse func(t *testing.T, res *pb.GetUserByEmailResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&repoResUser, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetUserByEmailResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, repoResUser.Username, res.User.Username)
				assert.Equal(t, repoResUser.Email, res.User.Email)
				assert.Equal(t, repoResUser.Password, res.User.Password)
				assert.NotNil(t, res.User.CreatedAt)
				assert.NotNil(t, res.User.UpdatedAt)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&domain.User{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.GetUserByEmailResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get user")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewUserUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewUserGRPCServer(server, usecase)
			res, err := s.GetUserByEmail(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestListUsers(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.ListUsersRequest
	}

	req := &pb.ListUsersRequest{
		Offset: 0,
		Limit:  10,
	}

	repoResUsers := []domain.User{
		{

			ID:        1,
			Username:  "test_username",
			Email:     "test@example.com",
			Password:  "test_password",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{

			ID:        2,
			Username:  "test_username",
			Email:     "test@example.com",
			Password:  "test_password",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, res *pb.ListUsersResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Return(&repoResUsers, nil)
			},
			checkResponse: func(t *testing.T, res *pb.ListUsersResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(repoResUsers), len(res.Users))
				for i, repoResUser := range repoResUsers {
					assert.Equal(t, repoResUser.Username, res.Users[i].Username)
					assert.Equal(t, repoResUser.Email, res.Users[i].Email)
					assert.Equal(t, repoResUser.Password, res.Users[i].Password)
					assert.NotNil(t, res.Users[i].CreatedAt)
					assert.NotNil(t, res.Users[i].UpdatedAt)
				}
			},
		},
		{
			name: "InvalidData",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Return(&[]domain.User{}, gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.ListUsersResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to list users")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewUserUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewUserGRPCServer(server, usecase)
			res, err := s.ListUsers(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateUserRequest
	}

	req := &pb.UpdateUserRequest{
		Id:       1,
		Username: "test_username",
		Email:    "test@example.com",
		Password: "test_password",
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().UpdateUser(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "InvalidArgument",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateUserRequest{},
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid argument")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewUserUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewUserGRPCServer(server, usecase)
			res, err := s.UpdateUser(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.DeleteUserRequest
	}

	req := &pb.DeleteUserRequest{
		Id: 1,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, res *pb.DeleteUserResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().DeleteUser(gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserResponse, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().DeleteUser(gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete user")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			usecase := usecase.NewUserUsecase(repo)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewUserGRPCServer(server, usecase)
			res, err := s.DeleteUser(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}
