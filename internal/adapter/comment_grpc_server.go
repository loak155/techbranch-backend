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

type ICommentGRPCServer interface {
	CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error)
	ListCommentsByUserID(ctx context.Context, req *pb.ListCommentsByUserIDRequest) (*pb.ListCommentsByUserIDResponse, error)
	ListCommentsByArticleID(ctx context.Context, req *pb.ListCommentsByArticleIDRequest) (*pb.ListCommentsByArticleIDResponse, error)
	DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error)
	DeleteCommentByUserIDAndArticleID(ctx context.Context, req *pb.DeleteCommentByUserIDAndArticleIDRequest) (*pb.DeleteCommentByUserIDAndArticleIDResponse, error)
	DeleteCommentByUserID(ctx context.Context, req *pb.DeleteCommentByUserIDRequest) (*pb.DeleteCommentByUserIDResponse, error)
	DeleteCommentByArticleID(ctx context.Context, req *pb.DeleteCommentByArticleIDRequest) (*pb.DeleteCommentByArticleIDResponse, error)
}

type commentGRPCServer struct {
	pb.UnimplementedCommentServiceServer
	usecase usecase.ICommentUsecase
}

func NewCommentGRPCServer(grpcServer *grpc.Server, usecase usecase.ICommentUsecase) pb.CommentServiceServer {
	server := commentGRPCServer{usecase: usecase}
	pb.RegisterCommentServiceServer(grpcServer, &server)
	return &server
}

func (server *commentGRPCServer) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	res := pb.CreateCommentResponse{}
	comment, err := server.usecase.CreateComment(
		domain.Comment{
			UserID:    uint(req.UserId),
			ArticleID: uint(req.ArticleId),
			Content:   req.Content,
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create comment: %v", err)
	}

	res.Comment = &pb.Comment{
		Id:        int32(comment.ID),
		UserId:    int32(comment.UserID),
		ArticleId: int32(comment.ArticleID),
		Content:   comment.Content,
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(comment.CreatedAt.Unix()), Nanos: int32(comment.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(comment.UpdatedAt.Unix()), Nanos: int32(comment.UpdatedAt.Nanosecond())},
	}

	return &res, nil
}

func (server *commentGRPCServer) ListCommentsByUserID(ctx context.Context, req *pb.ListCommentsByUserIDRequest) (*pb.ListCommentsByUserIDResponse, error) {
	res := pb.ListCommentsByUserIDResponse{}
	commentRes, err := server.usecase.ListCommentsByUserID(int(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list comments by user id: %v", err)
	}
	for _, comment := range commentRes {
		res.Comments = append(res.Comments, &pb.Comment{
			Id:        int32(comment.ID),
			UserId:    int32(comment.UserID),
			ArticleId: int32(comment.ArticleID),
			Content:   comment.Content,
			CreatedAt: &timestamppb.Timestamp{Seconds: int64(comment.CreatedAt.Unix()), Nanos: int32(comment.CreatedAt.Nanosecond())},
			UpdatedAt: &timestamppb.Timestamp{Seconds: int64(comment.UpdatedAt.Unix()), Nanos: int32(comment.UpdatedAt.Nanosecond())},
		})
	}

	return &res, nil
}

func (server *commentGRPCServer) ListCommentsByArticleID(ctx context.Context, req *pb.ListCommentsByArticleIDRequest) (*pb.ListCommentsByArticleIDResponse, error) {
	res := pb.ListCommentsByArticleIDResponse{}
	commentRes, err := server.usecase.ListCommentsByArticleID(int(req.ArticleId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list comments by article id: %v", err)
	}
	for _, comment := range commentRes {
		res.Comments = append(res.Comments, &pb.Comment{
			Id:        int32(comment.ID),
			UserId:    int32(comment.UserID),
			ArticleId: int32(comment.ArticleID),
			Content:   comment.Content,
			CreatedAt: &timestamppb.Timestamp{Seconds: int64(comment.CreatedAt.Unix()), Nanos: int32(comment.CreatedAt.Nanosecond())},
			UpdatedAt: &timestamppb.Timestamp{Seconds: int64(comment.UpdatedAt.Unix()), Nanos: int32(comment.UpdatedAt.Nanosecond())},
		})
	}

	return &res, nil
}

func (server *commentGRPCServer) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	res := pb.DeleteCommentResponse{}
	err := server.usecase.DeleteComment(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete comment: %v", err)
	}

	return &res, err
}

func (server *commentGRPCServer) DeleteCommentByUserIDAndArticleID(ctx context.Context, req *pb.DeleteCommentByUserIDAndArticleIDRequest) (*pb.DeleteCommentByUserIDAndArticleIDResponse, error) {
	res := pb.DeleteCommentByUserIDAndArticleIDResponse{}
	err := server.usecase.DeleteCommentByUserIDAndArticleID(int(req.UserId), int(req.ArticleId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete comment by user id and article id: %v", err)
	}

	return &res, err
}

func (server *commentGRPCServer) DeleteCommentByUserID(ctx context.Context, req *pb.DeleteCommentByUserIDRequest) (*pb.DeleteCommentByUserIDResponse, error) {
	res := pb.DeleteCommentByUserIDResponse{}
	err := server.usecase.DeleteCommentByUserID(int(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete comment by user id and article id: %v", err)
	}

	return &res, err
}

func (server *commentGRPCServer) DeleteCommentByArticleID(ctx context.Context, req *pb.DeleteCommentByArticleIDRequest) (*pb.DeleteCommentByArticleIDResponse, error) {
	res := pb.DeleteCommentByArticleIDResponse{}
	err := server.usecase.DeleteCommentByArticleID(int(req.ArticleId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete comment by user id and article id: %v", err)
	}

	return &res, err
}
