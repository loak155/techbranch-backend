package usecase

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/mock"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/pkg/mail"
	"github.com/loak155/techbranch-backend/pkg/oauth"
	"github.com/loak155/techbranch-backend/pkg/password"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPreSignup(t *testing.T) {
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
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			args: args{
				user: reqUser,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&domain.User{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "InvalidData",
			args: args{
				user: domain.User{},
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&domain.User{}, gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "UserAlreadyExists",
			args: args{
				user: domain.User{},
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&reqUser, gorm.ErrInvalidData)
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	testSMTPServer := smtpmock.New(smtpmock.ConfigurationAttr{})
	if err := testSMTPServer.Start(); err != nil {
		t.Fatal(err)
	}
	defer testSMTPServer.Stop()

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
			preSignupRedisManager := mock.NewRedisMock(t, 2, time.Duration(time.Hour*1))
			preSignupMailManager, err := mail.NewPresignupMailManager("localhost", testSMTPServer.PortNumber(), "test@example.com", "", "Test Prsignup", "../../pkg/mail/presignup.tmpl", "http://localhost:8080/v1/signup?token=")
			if err != nil {
				t.Fatalf("failed to create presignup mail manager: %v", err)
			}
			usecase := NewAuthUsecase(repo, *jwtAccessTokenManager, *jwtRefreshTokenManager, *redisAccessTokenManager, *redisRefreshTokenManager, *googleManager, *preSignupRedisManager, *preSignupMailManager)
			err = usecase.PreSignup(tc.args.user)
			tc.checkResponse(t, err)
		})
	}
}

func TestSignin(t *testing.T) {
	type args struct {
		email    string
		password string
	}

	reqEmail := "test@example.com"
	reqPassword := "test_password"

	hashedPassword, _ := password.HashPassword(reqPassword)

	repoResUser := domain.User{
		ID:        1,
		Username:  "test_username",
		Email:     "test@example.com",
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repoResUnmatchUser := domain.User{
		ID:        1,
		Username:  "test_username",
		Email:     "test@example.com",
		Password:  "unmatched_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, accessToken, refreshToken string, accessTokenExpiresIn, refreshTokenExpiresIn int, err error)
	}{
		{
			name: "OK",
			args: args{
				email:    reqEmail,
				password: reqPassword,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&repoResUser, nil)
			},
			checkResponse: func(t *testing.T, accessToken, refreshToken string, accessTokenExpiresIn, refreshTokenExpiresIn int, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, accessToken)
				assert.NotNil(t, refreshToken)
				assert.NotNil(t, accessTokenExpiresIn)
				assert.NotNil(t, refreshTokenExpiresIn)
			},
		},
		{
			name: "NotFound",
			args: args{
				email:    reqEmail,
				password: reqPassword,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&domain.User{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, accessToken, refreshToken string, accessTokenExpiresIn, refreshTokenExpiresIn int, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "UnmatchedPassword",
			args: args{
				email:    reqEmail,
				password: reqPassword,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&repoResUnmatchUser, nil)
			},
			checkResponse: func(t *testing.T, accessToken, refreshToken string, accessTokenExpiresIn, refreshTokenExpiresIn int, err error) {
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
			preSignupRedisManager := mock.NewRedisMock(t, 2, time.Duration(time.Hour*1))
			preSignupMailManager, _ := mail.NewPresignupMailManager("localhost", 2525, "test@example.com", "", "Test Prsignup", "../../pkg/mail/presignup.tmpl", "http://localhost:8080/v1/signup?token=")
			usecase := NewAuthUsecase(repo, *jwtAccessTokenManager, *jwtRefreshTokenManager, *redisAccessTokenManager, *redisRefreshTokenManager, *googleManager, *preSignupRedisManager, *preSignupMailManager)
			accessToken, refreshToken, accessTokenExpiresIn, refreshTokenExpiresIn, err := usecase.Signin(tc.args.email, tc.args.password)
			tc.checkResponse(t, accessToken, refreshToken, accessTokenExpiresIn, refreshTokenExpiresIn, err)
		})
	}
}

func TestSignout(t *testing.T) {
	type args struct {
		userID int
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
				userID: 1,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
			},
			checkResponse: func(t *testing.T, err error) {
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
			preSignupRedisManager := mock.NewRedisMock(t, 2, time.Duration(time.Hour*1))
			preSignupMailManager, _ := mail.NewPresignupMailManager("localhost", 2525, "test@example.com", "", "Test Prsignup", "../../pkg/mail/presignup.tmpl", "http://localhost:8080/v1/signup?token=")
			usecase := NewAuthUsecase(repo, *jwtAccessTokenManager, *jwtRefreshTokenManager, *redisAccessTokenManager, *redisRefreshTokenManager, *googleManager, *preSignupRedisManager, *preSignupMailManager)
			err := usecase.Signout(tc.args.userID)
			tc.checkResponse(t, err)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	type args struct {
		refreshToken string
	}

	jwtRefreshTokenManager := jwt.NewJwtManager("issuer", "secret", time.Duration(time.Hour*24*30))
	refreshToken, _, err := jwtRefreshTokenManager.GenerateToken(1)
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(repo *mock.MockIUserRepository)
		checkResponse func(t *testing.T, accessToken string, accessTokenExpiresIn int, err error)
	}{
		{
			name: "OK",
			args: args{
				refreshToken: refreshToken,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
			},
			checkResponse: func(t *testing.T, accessToken string, accessTokenExpiresIn int, err error) {
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
			preSignupRedisManager := mock.NewRedisMock(t, 2, time.Duration(time.Hour*1))
			preSignupMailManager, _ := mail.NewPresignupMailManager("localhost", 2525, "test@example.com", "", "Test Prsignup", "../../pkg/mail/presignup.tmpl", "http://localhost:8080/v1/signup?token=")
			usecase := NewAuthUsecase(repo, *jwtAccessTokenManager, *jwtRefreshTokenManager, *redisAccessTokenManager, *redisRefreshTokenManager, *googleManager, *preSignupRedisManager, *preSignupMailManager)
			accessToken, accessTokenExpiresIn, err := usecase.RefreshToken(tc.args.refreshToken)
			tc.checkResponse(t, accessToken, accessTokenExpiresIn, err)
		})
	}
}

func TestGetSigninUser(t *testing.T) {
	type args struct {
		userID int
	}

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
		checkResponse func(t *testing.T, user domain.User, err error)
	}{
		{
			name: "OK",
			args: args{
				userID: 1,
			},
			buildStubs: func(repo *mock.MockIUserRepository) {
				repo.EXPECT().GetUser(gomock.Any()).Return(&repoResUser, nil)
			},
			checkResponse: func(t *testing.T, user domain.User, err error) {
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
			preSignupRedisManager := mock.NewRedisMock(t, 2, time.Duration(time.Hour*1))
			preSignupMailManager, _ := mail.NewPresignupMailManager("localhost", 2525, "test@example.com", "", "Test Prsignup", "../../pkg/mail/presignup.tmpl", "http://localhost:8080/v1/signup?token=")
			usecase := NewAuthUsecase(repo, *jwtAccessTokenManager, *jwtRefreshTokenManager, *redisAccessTokenManager, *redisRefreshTokenManager, *googleManager, *preSignupRedisManager, *preSignupMailManager)
			user, err := usecase.GetSigninUser(tc.args.userID)
			tc.checkResponse(t, user, err)
		})
	}
}
