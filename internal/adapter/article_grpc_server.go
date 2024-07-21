package adapter

import (
	"context"

	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/internal/usecase"
	"github.com/loak155/techbranch-backend/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IArticleGRPCServer interface {
	CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error)
	GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleResponse, error)
	ListArticles(ctx context.Context, req *pb.ListArticlesRequest) (*pb.ListArticlesResponse, error)
	UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error)
	DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.DeleteArticleResponse, error)
}

type articleGRPCServer struct {
	pb.UnimplementedArticleServiceServer
	usecase usecase.IArticleUsecase
}

func NewArticleGRPCServer(grpcServer *grpc.Server, usecase usecase.IArticleUsecase) pb.ArticleServiceServer {
	server := articleGRPCServer{usecase: usecase}
	pb.RegisterArticleServiceServer(grpcServer, &server)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("article-server", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)
	return &server
}

func (server *articleGRPCServer) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	res := pb.CreateArticleResponse{}
	article, err := server.usecase.CreateArticle(
		domain.Article{
			Title: req.Title,
			Url:   req.Url,
			Image: req.Image,
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create article: %v", err)
	}

	res.Article = &pb.Article{
		Id:        int32(article.ID),
		Title:     article.Title,
		Url:       article.Url,
		Image:     article.Image,
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(article.CreatedAt.Unix()), Nanos: int32(article.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(article.UpdatedAt.Unix()), Nanos: int32(article.UpdatedAt.Nanosecond())},
	}

	return &res, nil
}

func (server *articleGRPCServer) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleResponse, error) {
	res := pb.GetArticleResponse{}
	article, err := server.usecase.GetArticle(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get article: %v", err)
	}
	res.Article = &pb.Article{
		Id:        int32(article.ID),
		Title:     article.Title,
		Url:       article.Url,
		Image:     article.Image,
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(article.CreatedAt.Unix()), Nanos: int32(article.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(article.UpdatedAt.Unix()), Nanos: int32(article.UpdatedAt.Nanosecond())},
	}

	return &res, nil
}

func (server *articleGRPCServer) ListArticles(ctx context.Context, req *pb.ListArticlesRequest) (*pb.ListArticlesResponse, error) {
	res := pb.ListArticlesResponse{}
	articleRes, err := server.usecase.ListArticles(int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list articles: %v", err)
	}
	for _, article := range articleRes {
		res.Articles = append(res.Articles, &pb.Article{
			Id:        int32(article.ID),
			Title:     article.Title,
			Url:       article.Url,
			Image:     article.Image,
			CreatedAt: &timestamppb.Timestamp{Seconds: int64(article.CreatedAt.Unix()), Nanos: int32(article.CreatedAt.Nanosecond())},
			UpdatedAt: &timestamppb.Timestamp{Seconds: int64(article.UpdatedAt.Unix()), Nanos: int32(article.UpdatedAt.Nanosecond())},
		})
	}

	return &res, nil
}

func (server *articleGRPCServer) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	res := pb.UpdateArticleResponse{}
	article, err := server.usecase.UpdateArticle(
		domain.Article{
			ID:    uint(req.Id),
			Title: req.Title,
			Url:   req.Url,
			Image: req.Image,
		},
	)

	res.Article = &pb.Article{
		Id:        int32(article.ID),
		Title:     article.Title,
		Url:       article.Url,
		Image:     article.Image,
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(article.CreatedAt.Unix()), Nanos: int32(article.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(article.UpdatedAt.Unix()), Nanos: int32(article.UpdatedAt.Nanosecond())},
	}
	return &res, err
}

func (server *articleGRPCServer) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.DeleteArticleResponse, error) {
	res := pb.DeleteArticleResponse{}
	err := server.usecase.DeleteArticle(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete article: %v", err)
	}

	return &res, err
}
