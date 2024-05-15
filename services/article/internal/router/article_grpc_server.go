package router

import (
	"context"

	"github.com/loak155/techbranch-backend/services/article/internal/domain"
	"github.com/loak155/techbranch-backend/services/article/internal/usecase"
	pb "github.com/loak155/techbranch-backend/services/article/proto"

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
	GetArticleCount(ctx context.Context, req *pb.GetArticleCountRequest) (*pb.GetArticleCountResponse, error)
	GetBookmarkArticle(ctx context.Context, req *pb.GetArticleCountRequest) (*pb.GetArticleCountResponse, error)
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
	article := domain.Article{Title: req.Title, Url: req.Url, Image: req.Image}
	articleRes, err := server.usecase.CreateArticle(ctx, article)
	if err != nil {
		return nil, err
	}
	res.Article = &pb.Article{
		Id:        int32(articleRes.ID),
		Title:     articleRes.Title,
		Url:       articleRes.Url,
		Image:     articleRes.Image,
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(articleRes.CreatedAt.Unix()), Nanos: int32(articleRes.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(articleRes.UpdatedAt.Unix()), Nanos: int32(articleRes.UpdatedAt.Nanosecond())},
	}

	return &res, nil
}

func (server *articleGRPCServer) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleResponse, error) {
	res := pb.GetArticleResponse{}
	articleRes, err := server.usecase.GetArticle(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	res.Article = &pb.Article{
		Id:        int32(articleRes.ID),
		Title:     articleRes.Title,
		Url:       articleRes.Url,
		Image:     articleRes.Image,
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(articleRes.CreatedAt.Unix()), Nanos: int32(articleRes.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(articleRes.UpdatedAt.Unix()), Nanos: int32(articleRes.UpdatedAt.Nanosecond())},
	}

	return &res, nil
}

func (server *articleGRPCServer) ListArticles(ctx context.Context, req *pb.ListArticlesRequest) (*pb.ListArticlesResponse, error) {
	res := pb.ListArticlesResponse{}
	articleRes, err := server.usecase.ListArticles(ctx, int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, err
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
	article := domain.Article{
		ID:        uint(req.Article.Id),
		Title:     req.Article.Title,
		Url:       req.Article.Url,
		Image:     req.Article.Image,
		CreatedAt: req.Article.CreatedAt.AsTime(),
		UpdatedAt: req.Article.UpdatedAt.AsTime(),
	}
	articleRes, err := server.usecase.UpdateArticle(ctx, article)
	res.Success = articleRes

	return &res, err
}

func (server *articleGRPCServer) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.DeleteArticleResponse, error) {
	res := pb.DeleteArticleResponse{}
	articleRes, err := server.usecase.DeleteArticle(ctx, int(req.Id))
	res.Success = articleRes

	return &res, err
}

func (server *articleGRPCServer) GetArticleCount(ctx context.Context, req *pb.GetArticleCountRequest) (*pb.GetArticleCountResponse, error) {
	res := pb.GetArticleCountResponse{}
	count, err := server.usecase.GetArticleCount(ctx)
	res.Count = int32(count)

	return &res, err
}

func (server *articleGRPCServer) GetBookmarkArticle(ctx context.Context, req *pb.GetBookmarkedArticleRequest) (*pb.GetBookmarkedArticleResponse, error) {
	res := pb.GetBookmarkedArticleResponse{}
	articles, err := server.usecase.GetBookmarkedArticle(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	for _, article := range articles {
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
