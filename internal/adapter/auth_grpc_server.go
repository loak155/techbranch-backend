package adapter

import (
	"context"

	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/internal/usecase"
	myContext "github.com/loak155/techbranch-backend/pkg/context"
	"github.com/loak155/techbranch-backend/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IAuthGRPCServer interface {
	PreSignup(ctx context.Context, req *pb.PreSignupRequest) (*pb.PreSignupResponse, error)
	Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error)
	Signin(ctx context.Context, req *pb.SigninRequest) (*pb.SigninResponse, error)
	Signout(ctx context.Context, req *pb.SignoutRequest) (*pb.SignoutResponse, error)
	RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error)
	GetSigninUser(ctx context.Context, req *pb.GetSigninUserRequest) (*pb.GetSigninUserResponse, error)
	GetGoogleLoginURL(ctx context.Context, req *pb.GetGoogleLoginURLRequest) (*pb.GetGoogleLoginURLResponse, error)
	GoogleLoginCallback(ctx context.Context, req *pb.GoogleLoginCallbackRequest) (*pb.GoogleLoginCallbackResponse, error)
}

type authGRPCServer struct {
	pb.UnimplementedAuthServiceServer
	usecase usecase.IAuthUsecase
}

func NewAuthGRPCServer(grpcServer *grpc.Server, usecase usecase.IAuthUsecase) pb.AuthServiceServer {
	server := authGRPCServer{usecase: usecase}
	pb.RegisterAuthServiceServer(grpcServer, &server)
	return &server
}

func (server *authGRPCServer) PreSignup(ctx context.Context, req *pb.PreSignupRequest) (*pb.PreSignupResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	err := server.usecase.PreSignup(
		domain.User{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to presignup: %v", err)
	}

	return &pb.PreSignupResponse{}, nil
}

func (server *authGRPCServer) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	if err := server.usecase.Signup(req.Token); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to signup: %v", err)
	}

	return &pb.SignupResponse{}, nil
}

func (server *authGRPCServer) Signin(ctx context.Context, req *pb.SigninRequest) (*pb.SigninResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	accessToken, refreshToken, expiresIn, err := server.usecase.Signin(req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to signin: %v", err)
	}
	res := pb.SigninResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int32(expiresIn),
		RefreshToken: refreshToken,
	}

	return &res, nil
}

func (server *authGRPCServer) Signout(ctx context.Context, req *pb.SignoutRequest) (*pb.SignoutResponse, error) {
	err := server.usecase.Signout(myContext.GetUserID(ctx))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to signout: %v", err)
	}

	return &pb.SignoutResponse{}, nil
}

func (server *authGRPCServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	accessToken, err := server.usecase.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to refresh token: %v", err)
	}
	res := pb.RefreshTokenResponse{
		AccessToken: accessToken,
	}

	return &res, nil
}

func (server *authGRPCServer) GetSigninUser(ctx context.Context, req *pb.GetSigninUserRequest) (*pb.GetSigninUserResponse, error) {
	res := pb.GetSigninUserResponse{}
	user, err := server.usecase.GetSigninUser(myContext.GetUserID(ctx))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get signin user: %v", err)
	}
	res.User = &pb.User{
		Id:        int32(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(user.CreatedAt.Unix()), Nanos: int32(user.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(user.UpdatedAt.Unix()), Nanos: int32(user.UpdatedAt.Nanosecond())},
	}

	return &res, nil
}

func (server *authGRPCServer) GetGoogleLoginURL(ctx context.Context, req *pb.GetGoogleLoginURLRequest) (*pb.GetGoogleLoginURLResponse, error) {
	res := pb.GetGoogleLoginURLResponse{}
	url := server.usecase.GetGoogleLoginURL()
	res.Url = url

	return &res, nil
}

func (server *authGRPCServer) GoogleLoginCallback(ctx context.Context, req *pb.GoogleLoginCallbackRequest) (*pb.GoogleLoginCallbackResponse, error) {
	accessToken, refreshToken, expiresIn, err := server.usecase.GoogleLoginCallback(req.State, req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to google login callback: %v", err)
	}
	res := pb.GoogleLoginCallbackResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int32(expiresIn),
		RefreshToken: refreshToken,
	}

	return &res, nil
}
