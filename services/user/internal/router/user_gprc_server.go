package router

import (
	"context"

	"github.com/loak155/techbranch-backend/services/user/internal/domain"
	"github.com/loak155/techbranch-backend/services/user/internal/usecase"
	pb "github.com/loak155/techbranch-backend/services/user/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IUserGRPCServer interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error)
	GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error)
	ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
}

type userGRPCServer struct {
	pb.UnimplementedUserServiceServer
	usecase usecase.IUserUsecase
}

func NewUserGRPCServer(grpcServer *grpc.Server, usecase usecase.IUserUsecase) pb.UserServiceServer {
	server := userGRPCServer{usecase: usecase}
	pb.RegisterUserServiceServer(grpcServer, &server)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("user-server", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)
	return &server
}

func (server *userGRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	res := pb.CreateUserResponse{}
	user := domain.User{Username: req.Username, Email: req.Email, Password: req.Password}
	userRes, err := server.usecase.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	res.User = &pb.User{
		Id:        int32(userRes.ID),
		Username:  userRes.Username,
		Email:     userRes.Email,
		Password:  userRes.Password,
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(userRes.CreatedAt.Unix()), Nanos: int32(userRes.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(userRes.UpdatedAt.Unix()), Nanos: int32(userRes.UpdatedAt.Nanosecond())},
	}

	return &res, nil
}

func (server *userGRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	res := pb.GetUserResponse{}
	userRes, err := server.usecase.GetUser(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	res.User = &pb.User{
		Id:        int32(userRes.ID),
		Username:  userRes.Username,
		Email:     userRes.Email,
		Password:  userRes.Password,
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(userRes.CreatedAt.Unix()), Nanos: int32(userRes.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(userRes.UpdatedAt.Unix()), Nanos: int32(userRes.UpdatedAt.Nanosecond())},
	}

	return &res, nil
}

func (server *userGRPCServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	res := pb.GetUserByEmailResponse{}
	userRes, err := server.usecase.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	res.User = &pb.User{
		Id:        int32(userRes.ID),
		Username:  userRes.Username,
		Email:     userRes.Email,
		Password:  userRes.Password,
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(userRes.CreatedAt.Unix()), Nanos: int32(userRes.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(userRes.UpdatedAt.Unix()), Nanos: int32(userRes.UpdatedAt.Nanosecond())},
	}

	return &res, nil
}

func (server *userGRPCServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	res := pb.ListUsersResponse{}
	userRes, err := server.usecase.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	for _, user := range userRes {
		res.Users = append(res.Users, &pb.User{
			Id:        int32(user.ID),
			Username:  user.Username,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: &timestamppb.Timestamp{Seconds: int64(user.CreatedAt.Unix()), Nanos: int32(user.CreatedAt.Nanosecond())},
			UpdatedAt: &timestamppb.Timestamp{Seconds: int64(user.UpdatedAt.Unix()), Nanos: int32(user.UpdatedAt.Nanosecond())},
		})
	}

	return &res, nil
}

func (server *userGRPCServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	res := pb.UpdateUserResponse{}
	user := domain.User{
		ID:        uint(req.User.Id),
		Username:  req.User.Username,
		Email:     req.User.Email,
		Password:  req.User.Password,
		CreatedAt: req.User.CreatedAt.AsTime(),
		UpdatedAt: req.User.UpdatedAt.AsTime(),
	}
	userRes, err := server.usecase.UpdateUser(ctx, user)
	res.Success = userRes

	return &res, err
}

func (server *userGRPCServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	res := pb.DeleteUserResponse{}
	userRes, err := server.usecase.DeleteUser(ctx, int(req.Id))
	res.Success = userRes

	return &res, err
}
