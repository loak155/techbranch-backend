package router

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	myContext "github.com/loak155/techbranch-backend/pkg/context"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/services/auth/internal/usecase"
	authpb "github.com/loak155/techbranch-backend/services/auth/proto"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type IAuthGRPCServer interface {
	Signup(ctx context.Context, req *authpb.SignupRequest) (*authpb.SignupResponse, error)
	Signin(ctx context.Context, req *authpb.SigninRequest) (*authpb.SigninResponse, error)
	GetSigninUser(ctx context.Context, req *authpb.GetSigninUserRequest) (*authpb.GetSigninUserResponse, error)
	GenerateToken(ctx context.Context, req *authpb.GenerateTokenRequest) (*authpb.GenerateTokenResponse, error)
	ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error)
	RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error)
	GetGoogleLoginURL(ctx context.Context, req *authpb.GetGoogleLoginURLRequest) (*authpb.GetGoogleLoginURLResponse, error)
	GoogleLoginCallback(ctx context.Context, req *authpb.GoogleLoginCallbackRequest) (authpb.GoogleLoginCallbackResponse, error)
}

type authGRPCServer struct {
	authpb.UnimplementedAuthServiceServer
	usecase    usecase.IAuthUsecase
	jwtManager jwt.JwtManager
}

func NewAuthGRPCServer(grpcServer *grpc.Server, usecase usecase.IAuthUsecase, jwtManager jwt.JwtManager) authpb.AuthServiceServer {
	server := authGRPCServer{usecase: usecase, jwtManager: jwtManager}
	authpb.RegisterAuthServiceServer(grpcServer, &server)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("auth-server", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)
	return &server
}

func (server *authGRPCServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	slog.Info("[Message]", "FullMethodName", fullMethodName)

	UnauthenticatedMethods := []string{
		"/techbranch.auth.AuthService/Signup",
		"/techbranch.auth.AuthService/Signin",
		"/techbranch.auth.AuthService/GetGoogleLoginURL",
		"/techbranch.auth.AuthService/GoogleLoginCallback",
	}

	// Allow some methods without authentication
	if slices.Contains(UnauthenticatedMethods, fullMethodName) {
		return ctx, nil
	}

	token, err := auth.AuthFromMD(ctx, "bearer")
	slog.Info("[Message]", "Token", token, "Error", err)
	if err != nil {
		return nil, err
	}

	claims, err := server.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	newCtx := myContext.SetUserID(ctx, claims.UserId)
	return newCtx, nil
}

func (server *authGRPCServer) Signup(ctx context.Context, req *authpb.SignupRequest) (*authpb.SignupResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	res := authpb.SignupResponse{}
	userRes, err := server.usecase.Signup(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	res.User = &authpb.User{
		Id: int32(userRes),
	}
	return &res, nil
}

func (server *authGRPCServer) Signin(ctx context.Context, req *authpb.SigninRequest) (*authpb.SigninResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	res := authpb.SigninResponse{}
	token, err := server.usecase.Signin(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	res.Token = token
	return &res, nil
}

func (server *authGRPCServer) GetSigninUser(ctx context.Context, req *authpb.GetSigninUserRequest) (*authpb.GetSigninUserResponse, error) {
	res := authpb.GetSigninUserResponse{}
	username, email, err := server.usecase.GetSigninUser(ctx, myContext.GetUserID(ctx))
	if err != nil {
		return nil, err
	}
	slog.Info("[Message]", "UserID", "確認2")
	res.User = &authpb.User{
		Id:       int32(myContext.GetUserID(ctx)),
		Username: username,
		Email:    email,
	}
	return &res, nil
}

func (server *authGRPCServer) GenerateToken(ctx context.Context, req *authpb.GenerateTokenRequest) (*authpb.GenerateTokenResponse, error) {
	res := authpb.GenerateTokenResponse{}
	authRes, err := server.usecase.GenerateToken(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	res.Token = authRes
	return &res, nil
}

func (server *authGRPCServer) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	res := authpb.ValidateTokenResponse{}
	authRes, err := server.usecase.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	res.Valid = authRes
	return &res, nil
}

func (server *authGRPCServer) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	res := authpb.RefreshTokenResponse{}
	authRes, err := server.usecase.RefreshToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	res.Token = authRes
	return &res, nil
}

func (server *authGRPCServer) GetGoogleLoginURL(ctx context.Context, req *authpb.GetGoogleLoginURLRequest) (*authpb.GetGoogleLoginURLResponse, error) {
	res := authpb.GetGoogleLoginURLResponse{}
	authRes, err := server.usecase.GetGoogleLoginURL(ctx)
	if err != nil {
		return nil, err
	}
	res.Url = authRes
	return &res, nil
}

func (server *authGRPCServer) GoogleLoginCallback(ctx context.Context, req *authpb.GoogleLoginCallbackRequest) (*authpb.GoogleLoginCallbackResponse, error) {
	res := authpb.GoogleLoginCallbackResponse{}
	authRes, err := server.usecase.GoogleLoginCallback(ctx, req.State, req.Code)
	if err != nil {
		return nil, err
	}
	res.Token = authRes
	return &res, nil
}
