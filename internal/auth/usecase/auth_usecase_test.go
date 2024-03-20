package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/loak155/techbranch-backend/internal/pkg/config"
	"github.com/loak155/techbranch-backend/internal/pkg/jwt"
	"github.com/loak155/techbranch-backend/internal/pkg/password"
	"github.com/loak155/techbranch-backend/mock"
	userpb "github.com/loak155/techbranch-backend/proto/user"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSignup(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
		email    string
		password string
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
		checkResponse func(t *testing.T, res int, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
				email:    "test@example.com",
				password: "password",
			},
			buildStubs: func(mockUserClient *mock.MockUserServiceClient) {
				mockUserClient.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(resUser, nil)
			},
			checkResponse: func(t *testing.T, res int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, res, int(resUser.User.Id))
			},
		},
		{
			name: "DuplicateEmail",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
				email:    "test@example.com",
				password: "password",
			},
			buildStubs: func(mockUserClient *mock.MockUserServiceClient) {
				mockUserClient.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(resUserErr, status.Error(codes.Internal, "failed to create user"))
			},
			checkResponse: func(t *testing.T, res int, err error) {
				assert.Error(t, err)
				assert.Equal(t, res, 0)
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

			usecase := NewAuthUsecase(mockUserClient, jwtManager)
			res, err := usecase.Signup(tc.args.ctx, tc.args.username, tc.args.email, tc.args.password)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestSignin(t *testing.T) {
	type args struct {
		ctx      context.Context
		email    string
		password string
	}

	hashedPassword, err := password.HashPassword("password")
	if err != nil {
		t.Errorf("failed to hash password: %v", err)
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
		checkResponse func(t *testing.T, res string, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "password",
			},
			buildStubs: func(mockUserClient *mock.MockUserServiceClient) {
				mockUserClient.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(resUser, nil)
			},
			checkResponse: func(t *testing.T, res string, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "password",
			},
			buildStubs: func(mockUserClient *mock.MockUserServiceClient) {
				mockUserClient.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(resUserErr, status.Error(codes.Internal, "failed to get user by email"))
			},
			checkResponse: func(t *testing.T, res string, err error) {
				assert.Error(t, err)
				assert.Equal(t, res, "")
			},
		},
		{
			name: "InvalidPassword",
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "InvalidPassword",
			},
			buildStubs: func(mockUserClient *mock.MockUserServiceClient) {
				mockUserClient.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(resUser, nil)
			},
			checkResponse: func(t *testing.T, res string, err error) {
				assert.Error(t, err)
				assert.Equal(t, res, "")
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

			usecase := NewAuthUsecase(mockUserClient, jwtManager)
			res, err := usecase.Signin(tc.args.ctx, tc.args.email, tc.args.password)
			tc.checkResponse(t, res, err)
		})
	}
}
