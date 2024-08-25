package adapter

import (
	"context"

	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/internal/usecase"
	"github.com/loak155/techbranch-backend/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IUserGRPCServer interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error)
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
	return &server
}

func (server *userGRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	res := pb.CreateUserResponse{}
	user, err := server.usecase.CreateUser(
		domain.User{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
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

func (server *userGRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	res := pb.GetUserResponse{}
	user, err := server.usecase.GetUser(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
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

func (server *userGRPCServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	res := pb.ListUsersResponse{}

	if req.Email != "" {
		if err := req.Validate(); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
		}
		user, err := server.usecase.GetUserByEmail(req.Email)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get user by email: %v", err)
		}
		res.Users = append(res.Users, &pb.User{
			Id:        int32(user.ID),
			Username:  user.Username,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: &timestamppb.Timestamp{Seconds: int64(user.CreatedAt.Unix()), Nanos: int32(user.CreatedAt.Nanosecond())},
			UpdatedAt: &timestamppb.Timestamp{Seconds: int64(user.UpdatedAt.Unix()), Nanos: int32(user.UpdatedAt.Nanosecond())},
		})
	} else {
		users, err := server.usecase.ListUsers(int(req.Offset), int(req.Limit))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
		}

		for _, user := range users {
			res.Users = append(res.Users, &pb.User{
				Id:        int32(user.ID),
				Username:  user.Username,
				Email:     user.Email,
				Password:  user.Password,
				CreatedAt: &timestamppb.Timestamp{Seconds: int64(user.CreatedAt.Unix()), Nanos: int32(user.CreatedAt.Nanosecond())},
				UpdatedAt: &timestamppb.Timestamp{Seconds: int64(user.UpdatedAt.Unix()), Nanos: int32(user.UpdatedAt.Nanosecond())},
			})
		}
		if len(users) == 0 {
			res.Users = append(res.Users, &pb.User{})
		}
	}

	return &res, nil
}

func (server *userGRPCServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	res := pb.UpdateUserResponse{}
	user, err := server.usecase.UpdateUser(
		domain.User{
			ID:       uint(req.Id),
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		},
	)

	res.User = &pb.User{
		Id:        int32(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(user.CreatedAt.Unix()), Nanos: int32(user.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(user.UpdatedAt.Unix()), Nanos: int32(user.UpdatedAt.Nanosecond())},
	}
	return &res, err
}

func (server *userGRPCServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	res := pb.DeleteUserResponse{}
	err := server.usecase.DeleteUser(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &res, err
}
