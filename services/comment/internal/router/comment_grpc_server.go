package router

import (
	"context"
	"log/slog"

	myContext "github.com/loak155/techbranch-backend/pkg/context"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/services/comment/internal/domain"
	"github.com/loak155/techbranch-backend/services/comment/internal/usecase"
	pb "github.com/loak155/techbranch-backend/services/comment/proto"
	"golang.org/x/exp/slices"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ICommentGRPCServer interface {
	CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error)
	GetComment(ctx context.Context, req *pb.GetCommentRequest) (*pb.GetCommentResponse, error)
	ListCommentsByArticleID(ctx context.Context, req *pb.ListCommentsByArticleIDRequest) (*pb.ListCommentsByArticleIDResponse, error)
	UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error)
	DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error)
	DeleteCommentByUserID(ctx context.Context, req *pb.DeleteCommentByUserIDRequest) (*pb.DeleteCommentByUserIDResponse, error)
	DeleteCommentByUserIDCompensate(ctx context.Context, req *pb.DeleteCommentByUserIDRequest) (*pb.DeleteCommentByUserIDResponse, error)
	DeleteCommentByArticleID(ctx context.Context, req *pb.DeleteCommentByArticleIDRequest) (*pb.DeleteCommentByArticleIDResponse, error)
	DeleteCommentByArticleIDCompensate(ctx context.Context, req *pb.DeleteCommentByArticleIDRequest) (*pb.DeleteCommentByArticleIDResponse, error)
}

type commentGRPCServer struct {
	pb.UnimplementedCommentServiceServer
	cu         usecase.ICommentUsecase
	jwtManager jwt.JwtManager
}

func NewCommentGRPCServer(grpcServer *grpc.Server, cu usecase.ICommentUsecase, jwtManager jwt.JwtManager) pb.CommentServiceServer {
	s := commentGRPCServer{cu: cu, jwtManager: jwtManager}
	pb.RegisterCommentServiceServer(grpcServer, &s)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("user-server", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)
	return &s
}

func (s *commentGRPCServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	slog.Info("[Message]", "FullMethodName", fullMethodName)

	UnauthenticatedMethods := []string{
		"/techbranch.comment.CommentService/ListCommentsByArticleID",
		"/techbranch.comment.CommentService/GetComment",
		"/techbranch.comment.CommentService/DeleteComment",
		"/techbranch.comment.CommentService/DeleteCommentByArticleID",
		"/techbranch.comment.CommentService/DeleteCommentByArticleIDCompensate",
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

	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	newCtx := myContext.SetUserID(ctx, claims.UserId)
	return newCtx, nil
}

func (s *commentGRPCServer) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	res := pb.CreateCommentResponse{}
	userId := myContext.GetUserID(ctx)

	commentRes, err := s.cu.CreateComment(
		ctx,
		domain.Comment{
			UserID:    uint(userId),
			ArticleID: uint(req.ArticleId),
			Content:   req.Content,
		},
	)
	if err != nil {
		return nil, err
	}
	res.Comment = &pb.Comment{
		Id:        int32(commentRes.ID),
		UserId:    int32(commentRes.UserID),
		ArticleId: int32(commentRes.ArticleID),
		Content:   commentRes.Content,
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(commentRes.CreatedAt.Unix()), Nanos: int32(commentRes.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(commentRes.UpdatedAt.Unix()), Nanos: int32(commentRes.UpdatedAt.Nanosecond())},
	}

	return &res, nil
}

func (s *commentGRPCServer) GetComment(ctx context.Context, req *pb.GetCommentRequest) (*pb.GetCommentResponse, error) {
	res := pb.GetCommentResponse{}
	commentRes, err := s.cu.GetComment(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	res.Comment = &pb.Comment{
		Id:        int32(commentRes.ID),
		UserId:    int32(commentRes.UserID),
		ArticleId: int32(commentRes.ArticleID),
		Content:   commentRes.Content,
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(commentRes.CreatedAt.Unix()), Nanos: int32(commentRes.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(commentRes.UpdatedAt.Unix()), Nanos: int32(commentRes.UpdatedAt.Nanosecond())},
	}

	return &res, nil
}

func (s *commentGRPCServer) ListCommentsByArticleID(ctx context.Context, req *pb.ListCommentsByArticleIDRequest) (*pb.ListCommentsByArticleIDResponse, error) {
	res := pb.ListCommentsByArticleIDResponse{}
	commentRes, err := s.cu.ListCommentsByArticleID(ctx, int(req.ArticleId))
	if err != nil {
		return nil, err
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

func (s *commentGRPCServer) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
	res := pb.UpdateCommentResponse{}
	comment := domain.Comment{
		ID:        uint(req.Comment.Id),
		UserID:    uint(req.Comment.UserId),
		ArticleID: uint(req.Comment.ArticleId),
		Content:   req.Comment.Content,
		CreatedAt: req.Comment.CreatedAt.AsTime(),
		UpdatedAt: req.Comment.UpdatedAt.AsTime(),
	}
	commentRes, err := s.cu.UpdateComment(ctx, comment)
	res.Success = commentRes

	return &res, err
}

func (s *commentGRPCServer) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	res := pb.DeleteCommentResponse{}
	commentRes, err := s.cu.DeleteComment(ctx, int(req.Id))
	res.Success = commentRes

	return &res, err
}

func (s *commentGRPCServer) DeleteCommentByUserID(ctx context.Context, req *pb.DeleteCommentByUserIDRequest) (*pb.DeleteCommentByUserIDResponse, error) {
	res := pb.DeleteCommentByUserIDResponse{}
	commentRes, err := s.cu.DeleteCommentByUserID(ctx, int(req.UserId))
	res.Success = commentRes

	return &res, err
}

func (s *commentGRPCServer) DeleteCommentByUserIDCompensate(ctx context.Context, req *pb.DeleteCommentByUserIDRequest) (*pb.DeleteCommentByUserIDResponse, error) {
	res := pb.DeleteCommentByUserIDResponse{}
	commentRes, err := s.cu.DeleteCommentByUserIDCompensate(ctx, int(req.UserId))
	res.Success = commentRes

	return &res, err
}

func (s *commentGRPCServer) DeleteCommentByArticleID(ctx context.Context, req *pb.DeleteCommentByArticleIDRequest) (*pb.DeleteCommentByArticleIDResponse, error) {
	res := pb.DeleteCommentByArticleIDResponse{}
	commentRes, err := s.cu.DeleteCommentByArticleID(ctx, int(req.ArticleId))
	res.Success = commentRes

	return &res, err
}

func (s *commentGRPCServer) DeleteCommentByArticleIDCompensate(ctx context.Context, req *pb.DeleteCommentByArticleIDRequest) (*pb.DeleteCommentByArticleIDResponse, error) {
	res := pb.DeleteCommentByArticleIDResponse{}
	commentRes, err := s.cu.DeleteCommentByArticleIDCompensate(ctx, int(req.ArticleId))
	res.Success = commentRes

	return &res, err
}
