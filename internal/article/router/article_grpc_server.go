package router

import (
	"context"

	"github.com/loak155/techbranch-backend/internal/article/domain"
	"github.com/loak155/techbranch-backend/internal/article/usecase"
	pb "github.com/loak155/techbranch-backend/proto/article"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IArticleGRPCServer interface {
	CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error)
	GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleResponse, error)
	ListArticles(ctx context.Context, req *pb.ListArticlesRequest) (*pb.ListArticlesResponse, error)
	UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error)
	DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.DeleteArticleResponse, error)
	IncrementBookmarksCount(ctx context.Context, req *pb.IncrementBookmarksCountRequest) (*pb.IncrementBookmarksCountResponse, error)
	IncrementBookmarksCountCompensate(ctx context.Context, req *pb.IncrementBookmarksCountRequest) (*pb.IncrementBookmarksCountResponse, error)
	DecrementBookmarksCount(ctx context.Context, req *pb.DecrementBookmarksCountRequest) (*pb.DecrementBookmarksCountResponse, error)
	DecrementBookmarksCountCompensate(ctx context.Context, req *pb.DecrementBookmarksCountRequest) (*pb.DecrementBookmarksCountResponse, error)
}

type articleGRPCServer struct {
	pb.UnimplementedArticleServiceServer
	usecase usecase.IArticleUsecase
}

func NewArticleGRPCServer(grpcServer *grpc.Server, usecase usecase.IArticleUsecase) pb.ArticleServiceServer {
	server := articleGRPCServer{usecase: usecase}
	pb.RegisterArticleServiceServer(grpcServer, &server)
	reflection.Register(grpcServer)
	return &server
}

func (server *articleGRPCServer) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error) {
	res := pb.CreateArticleResponse{}
	article := domain.Article{Title: req.Article.Title, Url: req.Article.Url}
	articleRes, err := server.usecase.CreateArticle(ctx, article)
	if err != nil {
		return nil, err
	}
	res.Article = &pb.Article{
		Id:            int32(articleRes.ID),
		Title:         articleRes.Title,
		Url:           articleRes.Url,
		BookmarkCount: int32(article.BookmarkCount),
		CreatedAt:     &timestamppb.Timestamp{Seconds: int64(articleRes.CreatedAt.Unix()), Nanos: int32(articleRes.CreatedAt.Nanosecond())},
		UpdatedAt:     &timestamppb.Timestamp{Seconds: int64(articleRes.UpdatedAt.Unix()), Nanos: int32(articleRes.UpdatedAt.Nanosecond())},
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
		Id:            int32(articleRes.ID),
		Title:         articleRes.Title,
		Url:           articleRes.Url,
		BookmarkCount: int32(articleRes.BookmarkCount),
		CreatedAt:     &timestamppb.Timestamp{Seconds: int64(articleRes.CreatedAt.Unix()), Nanos: int32(articleRes.CreatedAt.Nanosecond())},
		UpdatedAt:     &timestamppb.Timestamp{Seconds: int64(articleRes.UpdatedAt.Unix()), Nanos: int32(articleRes.UpdatedAt.Nanosecond())},
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
			Id:            int32(article.ID),
			Title:         article.Title,
			Url:           article.Url,
			BookmarkCount: int32(article.BookmarkCount),
			CreatedAt:     &timestamppb.Timestamp{Seconds: int64(article.CreatedAt.Unix()), Nanos: int32(article.CreatedAt.Nanosecond())},
			UpdatedAt:     &timestamppb.Timestamp{Seconds: int64(article.UpdatedAt.Unix()), Nanos: int32(article.UpdatedAt.Nanosecond())},
		})
	}

	return &res, nil
}

func (server *articleGRPCServer) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error) {
	res := pb.UpdateArticleResponse{}
	article := domain.Article{
		ID:            uint(req.Article.Id),
		Title:         req.Article.Title,
		Url:           req.Article.Url,
		BookmarkCount: uint(req.Article.BookmarkCount),
		CreatedAt:     req.Article.CreatedAt.AsTime(),
		UpdatedAt:     req.Article.UpdatedAt.AsTime(),
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

func (server *articleGRPCServer) IncrementBookmarksCount(ctx context.Context, req *pb.IncrementBookmarksCountRequest) (*pb.IncrementBookmarksCountResponse, error) {
	res, err := server.usecase.IncrementBookmarksCount(ctx, int(req.Id))
	return &pb.IncrementBookmarksCountResponse{Success: res}, err
}

func (server *articleGRPCServer) IncrementBookmarksCountCompensate(ctx context.Context, req *pb.IncrementBookmarksCountRequest) (*pb.IncrementBookmarksCountResponse, error) {
	res, err := server.usecase.IncrementBookmarksCountCompensate(ctx, int(req.Id))
	return &pb.IncrementBookmarksCountResponse{Success: res}, err
}

func (server *articleGRPCServer) DecrementBookmarksCount(ctx context.Context, req *pb.DecrementBookmarksCountRequest) (*pb.DecrementBookmarksCountResponse, error) {
	res, err := server.usecase.DecrementBookmarksCount(ctx, int(req.Id))
	return &pb.DecrementBookmarksCountResponse{Success: res}, err
}

func (server *articleGRPCServer) DecrementBookmarksCountCompensate(ctx context.Context, req *pb.DecrementBookmarksCountRequest) (*pb.DecrementBookmarksCountResponse, error) {
	res, err := server.usecase.DecrementBookmarksCountCompensate(ctx, int(req.Id))
	return &pb.DecrementBookmarksCountResponse{Success: res}, err
}
