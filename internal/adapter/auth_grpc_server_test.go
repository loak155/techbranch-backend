package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/internal/usecase"
	"github.com/loak155/techbranch-backend/mock"
	myContext "github.com/loak155/techbranch-backend/pkg/context"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/pkg/oauth"
	"github.com/loak155/techbranch-backend/pkg/password"
	"github.com/loak155/techbranch-backend/pkg/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func TestSignup(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.SignupRequest
	}

	req := &pb.SignupRequest{
		Username: "test_user",
		Email:    "test@example.com",
		Password: "password",
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, res *pb.SignupResponse, err error)
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
			checkResponse: func(t *testing.T, res *pb.SignupResponse, err error) {
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
				req: &pb.SignupRequest{},
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
			},
			checkResponse: func(t *testing.T, res *pb.SignupResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid argument")
			},
		},
		{
			name: "FailedToCreateUser",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any()).Return(gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.SignupResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to signup")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			jwtAccessTokenManager := jwt.NewJwtManager("issuer", "secret", time.Duration(time.Hour*1))
			jwtRefreshTokenManager := jwt.NewJwtManager("issuer", "secret", time.Duration(time.Hour*24*30))
			redisAccessTokenManager := mock.NewRedisMock(t, 0, time.Duration(time.Hour*1))
			redisRefreshTokenManager := mock.NewRedisMock(t, 1, time.Duration(time.Hour*24*30))
			googleManager := oauth.NewGoogleManager("state", "clientID", "clientSecret", "redirectURL")
			usecase := usecase.NewAuthUsecase(repo, *jwtAccessTokenManager, *jwtRefreshTokenManager, *redisAccessTokenManager, *redisRefreshTokenManager, *googleManager)

			server := grpc.NewServer()
			server.GracefulStop()

			s := NewAuthGRPCServer(server, usecase)
			res, err := s.Signup(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestSignin(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.SigninRequest
	}

	req := &pb.SigninRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	hashedPassword, _ := password.HashPassword("password")
	repoResUser := domain.User{
		ID:        1,
		Username:  "test_username",
		Email:     "test@example.com",
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, res *pb.SigninResponse, err error)
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
			checkResponse: func(t *testing.T, res *pb.SigninResponse, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, res.AccessToken)
				assert.Equal(t, "Bearer", res.TokenType)
				assert.NotEmpty(t, res.RefreshToken)
				assert.NotEmpty(t, res.ExpiresIn)
			},
		},
		{
			name: "InvalidArgument",
			args: args{
				ctx: context.Background(),
				req: &pb.SigninRequest{},
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
			},
			checkResponse: func(t *testing.T, res *pb.SigninResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid argument")
			},
		},
		{
			name: "FailedToSignin",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&domain.User{}, gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.SigninResponse, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to signin")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			jwtAccessTokenManager := jwt.NewJwtManager("issuer", "secret", time.Duration(time.Hour*1))
			jwtRefreshTokenManager := jwt.NewJwtManager("issuer", "secret", time.Duration(time.Hour*24*30))
			redisAccessTokenManager := mock.NewRedisMock(t, 0, time.Duration(time.Hour*1))
			redisRefreshTokenManager := mock.NewRedisMock(t, 1, time.Duration(time.Hour*24*30))
			googleManager := oauth.NewGoogleManager("state", "clientID", "clientSecret", "redirectURL")
			usecase := usecase.NewAuthUsecase(repo, *jwtAccessTokenManager, *jwtRefreshTokenManager, *redisAccessTokenManager, *redisRefreshTokenManager, *googleManager)

			server := grpc.NewServer()
			server.GracefulStop()

			s := NewAuthGRPCServer(server, usecase)
			res, err := s.Signin(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestSignout(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.SignoutRequest
	}

	ctx := myContext.SetUserID(context.Background(), 1)
	req := &pb.SignoutRequest{}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, res *pb.SignoutResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: ctx,
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
			},
			checkResponse: func(t *testing.T, res *pb.SignoutResponse, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mock.NewMockIUserRepository(mockCtrl)
			tc.buildStubs(repo)

			jwtAccessTokenManager := jwt.NewJwtManager("issuer", "secret", time.Duration(time.Hour*1))
			jwtRefreshTokenManager := jwt.NewJwtManager("issuer", "secret", time.Duration(time.Hour*24*30))
			redisAccessTokenManager := mock.NewRedisMock(t, 0, time.Duration(time.Hour*1))
			redisRefreshTokenManager := mock.NewRedisMock(t, 1, time.Duration(time.Hour*24*30))
			googleManager := oauth.NewGoogleManager("state", "clientID", "clientSecret", "redirectURL")
			usecase := usecase.NewAuthUsecase(repo, *jwtAccessTokenManager, *jwtRefreshTokenManager, *redisAccessTokenManager, *redisRefreshTokenManager, *googleManager)

			server := grpc.NewServer()
			server.GracefulStop()

			s := NewAuthGRPCServer(server, usecase)
			res, err := s.Signout(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.RefreshTokenRequest
	}

	jwtRefreshTokenManager := jwt.NewJwtManager("issuer", "secret", time.Duration(time.Hour*24*30))
	refreshToken, _, _ := jwtRefreshTokenManager.GenerateToken(1)
	req := &pb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, res *pb.RefreshTokenResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
			},
			checkResponse: func(t *testing.T, res *pb.RefreshTokenResponse, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "InvalidToken",
			args: args{
				ctx: context.Background(),
				req: &pb.RefreshTokenRequest{},
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
			},
			checkResponse: func(t *testing.T, res *pb.RefreshTokenResponse, err error) {
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

			jwtAccessTokenManager := jwt.NewJwtManager("issuer", "secret", time.Duration(time.Hour*1))
			jwtRefreshTokenManager := jwt.NewJwtManager("issuer", "secret", time.Duration(time.Hour*24*30))
			redisAccessTokenManager := mock.NewRedisMock(t, 0, time.Duration(time.Hour*1))
			redisRefreshTokenManager := mock.NewRedisMock(t, 1, time.Duration(time.Hour*24*30))
			googleManager := oauth.NewGoogleManager("state", "clientID", "clientSecret", "redirectURL")
			usecase := usecase.NewAuthUsecase(repo, *jwtAccessTokenManager, *jwtRefreshTokenManager, *redisAccessTokenManager, *redisRefreshTokenManager, *googleManager)

			server := grpc.NewServer()
			server.GracefulStop()

			s := NewAuthGRPCServer(server, usecase)
			res, err := s.RefreshToken(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestGetSigninUser(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.GetSigninUserRequest
	}

	ctx := myContext.SetUserID(context.Background(), 1)
	req := &pb.GetSigninUserRequest{}

	repoResUser := domain.User{
		ID:        1,
		Username:  "test_username",
		Email:     "test@example.com",
		Password:  "hashedPassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, res *pb.GetSigninUserResponse, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: ctx,
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUser(gomock.Any()).Return(&repoResUser, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetSigninUserResponse, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "FailedToGetUser",
			args: args{
				ctx: ctx,
				req: req,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUser(gomock.Any()).Return(&domain.User{}, gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, res *pb.GetSigninUserResponse, err error) {
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

			jwtAccessTokenManager := jwt.NewJwtManager("issuer", "secret", time.Duration(time.Hour*1))
			jwtRefreshTokenManager := jwt.NewJwtManager("issuer", "secret", time.Duration(time.Hour*24*30))
			redisAccessTokenManager := mock.NewRedisMock(t, 0, time.Duration(time.Hour*1))
			redisRefreshTokenManager := mock.NewRedisMock(t, 1, time.Duration(time.Hour*24*30))
			googleManager := oauth.NewGoogleManager("state", "clientID", "clientSecret", "redirectURL")
			usecase := usecase.NewAuthUsecase(repo, *jwtAccessTokenManager, *jwtRefreshTokenManager, *redisAccessTokenManager, *redisRefreshTokenManager, *googleManager)

			server := grpc.NewServer()
			server.GracefulStop()

			s := NewAuthGRPCServer(server, usecase)
			res, err := s.GetSigninUser(tc.args.ctx, tc.args.req)
			tc.checkResponse(t, res, err)
		})
	}
}
