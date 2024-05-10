package router

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/loak155/techbranch-backend/pkg/config"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/pkg/oauth"
	"github.com/loak155/techbranch-backend/pkg/password"
	"github.com/loak155/techbranch-backend/services/auth/internal/usecase"
	authpb "github.com/loak155/techbranch-backend/services/auth/proto"
	"github.com/loak155/techbranch-backend/services/user/mock"
	userpb "github.com/loak155/techbranch-backend/services/user/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSignup(t *testing.T) {
	type args struct {
		ctx context.Context
		req *authpb.SignupRequest
	}

	req := &authpb.SignupRequest{
		Username: "test_user",
		Email:    "test@example.com",
		Password: "password",
	}

	resUser := &userpb.CreateUserResponse{
		User: &userpb.User{
			Id:       1,
			Username: "test_user",
			Email:    "test@example.com",
			Password: "password",
		},
	}

	resUserErr := &userpb.CreateUserResponse{
		User: &userpb.User{},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(mockUserClient *mock.MockUserServiceClient)
		checkResponse func(t *testing.T, res *authpb.SignupResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(mockUserClient *mock.MockUserServiceClient) {
				mockUserClient.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(resUser, nil)
			},
			checkResponse: func(t *testing.T, res *authpb.SignupResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, int(res.User.Id))
			},
		},
		{
			name: "DuplicateEmail",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(mockUserClient *mock.MockUserServiceClient) {
				mockUserClient.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(resUserErr, status.Error(codes.Internal, "failed to create user"))
			},
			checkResponse: func(t *testing.T, res *authpb.SignupResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf, err := config.Load("../../../configs/config.yaml")
			assert.NoError(t, err)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockUserClient := mock.NewMockUserServiceClient(mockCtrl)
			tc.buildStubs(mockUserClient)
			jwtManager := jwt.NewJwtManager(conf)
			google := oauth.NewGoogle(conf)

			usecase := usecase.NewAuthUsecase(mockUserClient, jwtManager, google)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewAuthGRPCServer(server, usecase, jwtManager)
			res, err := s.Signup(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestSignin(t *testing.T) {
	type args struct {
		ctx context.Context
		req *authpb.SigninRequest
	}

	hashedPassword, err := password.HashPassword("password")
	if err != nil {
		t.Errorf("failed to hash password: %v", err)
	}

	req := &authpb.SigninRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	reqInvalidPassword := &authpb.SigninRequest{
		Email:    "test@example.com",
		Password: "InvalidPassword",
	}

	resUser := &userpb.GetUserByEmailResponse{
		User: &userpb.User{
			Id:       1,
			Username: "test_user",
			Email:    "test@example.com",
			Password: hashedPassword,
		},
	}

	resUserErr := &userpb.GetUserByEmailResponse{
		User: &userpb.User{},
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(mockUserClient *mock.MockUserServiceClient)
		checkResponse func(t *testing.T, res *authpb.SigninResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(mockUserClient *mock.MockUserServiceClient) {
				mockUserClient.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(resUser, nil)
			},
			checkResponse: func(t *testing.T, res *authpb.SigninResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(mockUserClient *mock.MockUserServiceClient) {
				mockUserClient.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(resUserErr, status.Error(codes.Internal, "failed to get user by email"))
			},
			checkResponse: func(t *testing.T, res *authpb.SigninResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "InvalidPassword",
			args: args{
				ctx: context.Background(),
				req: reqInvalidPassword,
			},
			buildStubs: func(mockUserClient *mock.MockUserServiceClient) {
				mockUserClient.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(resUser, nil)
			},
			checkResponse: func(t *testing.T, res *authpb.SigninResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf, err := config.Load("../../../configs/config.yaml")
			assert.NoError(t, err)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockUserClient := mock.NewMockUserServiceClient(mockCtrl)
			tc.buildStubs(mockUserClient)
			jwtManager := jwt.NewJwtManager(conf)
			google := oauth.NewGoogle(conf)

			usecase := usecase.NewAuthUsecase(mockUserClient, jwtManager, google)
			server := grpc.NewServer()
			server.GracefulStop()

			s := NewAuthGRPCServer(server, usecase, jwtManager)
			res, err := s.Signin(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}
