package router

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	myContext "github.com/loak155/techbranch-backend/pkg/context"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/services/bookmark/internal/domain"
	"github.com/loak155/techbranch-backend/services/bookmark/internal/usecase"
	pb "github.com/loak155/techbranch-backend/services/bookmark/proto"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IBookmarkGRPCServer interface {
	CreateBookmark(ctx context.Context, req *pb.CreateBookmarkRequest) (*pb.CreateBookmarkResponse, error)
	GetBookmark(context.Context, *pb.GetBookmarkRequest) (*pb.GetBookmarkResponse, error)
	GetBookmarkCountByArticleID(context.Context, *pb.GetBookmarkCountByArticleIDRequest) (*pb.GetBookmarkCountByArticleIDResponse, error)
	ListBookmarks(context.Context, *pb.ListBookmarksRequest) (*pb.ListBookmarksResponse, error)
	ListBookmarksByUserID(context.Context, *pb.ListBookmarksByUserIDRequest) (*pb.ListBookmarksByUserIDResponse, error)
	ListBookmarksByArticleID(context.Context, *pb.ListBookmarksByArticleIDRequest) (*pb.ListBookmarksByArticleIDResponse, error)
	DeleteBookmark(ctx context.Context, req *pb.DeleteBookmarkRequest) (*pb.DeleteBookmarkResponse, error)
	DeleteBookmarkByUserID(ctx context.Context, req *pb.DeleteBookmarkByUserIDRequest) (*pb.DeleteBookmarkByUserIDResponse, error)
	DeleteBookmarkByUserIDCompensate(ctx context.Context, req *pb.DeleteBookmarkByUserIDRequest) (*pb.DeleteBookmarkByUserIDResponse, error)
	DeleteBookmarkByArticleID(ctx context.Context, req *pb.DeleteBookmarkByArticleIDRequest) (*pb.DeleteBookmarkByArticleIDResponse, error)
	DeleteBookmarkByArticleIDCompensate(ctx context.Context, req *pb.DeleteBookmarkByArticleIDRequest) (*pb.DeleteBookmarkByArticleIDResponse, error)
}

type bookmarkGRPCServer struct {
	pb.UnimplementedBookmarkServiceServer
	usecase    usecase.IBookmarkUsecase
	jwtManager jwt.JwtManager
}

func NewBookmarkGRPCServer(grpcServer *grpc.Server, usecase usecase.IBookmarkUsecase, jwtManager jwt.JwtManager) pb.BookmarkServiceServer {
	s := bookmarkGRPCServer{usecase: usecase, jwtManager: jwtManager}
	pb.RegisterBookmarkServiceServer(grpcServer, &s)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("user-server", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)
	return &s
}

func (server *bookmarkGRPCServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	slog.Info("[Message]", "FullMethodName", fullMethodName)

	UnauthenticatedMethods := []string{
		"/techbranch.bookmark.BookmarkService/GetBookmark",
		"/techbranch.bookmark.BookmarkService/GetBookmarkCountByArticleID",
		"/techbranch.bookmark.BookmarkService/ListBookmarks",
		"/techbranch.bookmark.BookmarkService/ListBookmarksByUserID",
		"/techbranch.bookmark.BookmarkService/ListBookmarksByArticleID",
		"/techbranch.bookmark.BookmarkService/DeleteBookmarkByUserID",
		"/techbranch.bookmark.BookmarkService/DeleteBookmarkByUserIDCompensate",
		"/techbranch.bookmark.BookmarkService/DeleteBookmarkByArticleID",
		"/techbranch.bookmark.BookmarkService/DeleteBookmarkByArticleIDCompensate",
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

func (server *bookmarkGRPCServer) CreateBookmark(ctx context.Context, req *pb.CreateBookmarkRequest) (*pb.CreateBookmarkResponse, error) {
	res := pb.CreateBookmarkResponse{}

	bookmarkRes, err := server.usecase.CreateBookmark(
		ctx,
		domain.Bookmark{
			UserID:    uint(req.UserId),
			ArticleID: uint(req.ArticleId),
		},
	)
	if err != nil {
		return nil, err
	}
	res.Bookmark = &pb.Bookmark{
		Id:        int32(bookmarkRes.ID),
		UserId:    int32(bookmarkRes.UserID),
		ArticleId: int32(bookmarkRes.ArticleID),
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(bookmarkRes.CreatedAt.Unix()), Nanos: int32(bookmarkRes.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(bookmarkRes.UpdatedAt.Unix()), Nanos: int32(bookmarkRes.UpdatedAt.Nanosecond())},
	}

	return &res, nil
}

func (server *bookmarkGRPCServer) GetBookmark(ctx context.Context, req *pb.GetBookmarkRequest) (*pb.GetBookmarkResponse, error) {
	res := pb.GetBookmarkResponse{}
	bookmarkRes, err := server.usecase.GetBookmark(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	res.Bookmark = &pb.Bookmark{
		Id:        int32(bookmarkRes.ID),
		UserId:    int32(bookmarkRes.UserID),
		ArticleId: int32(bookmarkRes.ArticleID),
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(bookmarkRes.CreatedAt.Unix()), Nanos: int32(bookmarkRes.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(bookmarkRes.UpdatedAt.Unix()), Nanos: int32(bookmarkRes.UpdatedAt.Nanosecond())},
	}

	return &res, nil
}

func (server *bookmarkGRPCServer) GetBookmarkCountByArticleID(ctx context.Context, req *pb.GetBookmarkCountByArticleIDRequest) (*pb.GetBookmarkCountByArticleIDResponse, error) {
	res := pb.GetBookmarkCountByArticleIDResponse{}
	count, err := server.usecase.GetBookmarkCountByArticleID(ctx, int(req.ArticleId))
	if err != nil {
		return nil, err
	}
	res.Count = int32(count)

	return &res, nil
}

func (server *bookmarkGRPCServer) ListBookmarks(ctx context.Context, req *pb.ListBookmarksRequest) (*pb.ListBookmarksResponse, error) {
	res := pb.ListBookmarksResponse{}
	bookmarkRes, err := server.usecase.ListBookmarks(ctx)
	if err != nil {
		return nil, err
	}
	for _, bookmark := range bookmarkRes {
		res.Bookmarks = append(res.Bookmarks, &pb.Bookmark{
			Id:        int32(bookmark.ID),
			UserId:    int32(bookmark.UserID),
			ArticleId: int32(bookmark.ArticleID),
			CreatedAt: &timestamppb.Timestamp{Seconds: int64(bookmark.CreatedAt.Unix()), Nanos: int32(bookmark.CreatedAt.Nanosecond())},
			UpdatedAt: &timestamppb.Timestamp{Seconds: int64(bookmark.UpdatedAt.Unix()), Nanos: int32(bookmark.UpdatedAt.Nanosecond())},
		})
	}

	return &res, nil
}

func (server *bookmarkGRPCServer) ListBookmarksByUserID(ctx context.Context, req *pb.ListBookmarksByUserIDRequest) (*pb.ListBookmarksByUserIDResponse, error) {
	res := pb.ListBookmarksByUserIDResponse{}
	bookmarkRes, err := server.usecase.ListBookmarksByUserID(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	for _, bookmark := range bookmarkRes {
		res.Bookmarks = append(res.Bookmarks, &pb.Bookmark{
			Id:        int32(bookmark.ID),
			UserId:    int32(bookmark.UserID),
			ArticleId: int32(bookmark.ArticleID),
			CreatedAt: &timestamppb.Timestamp{Seconds: int64(bookmark.CreatedAt.Unix()), Nanos: int32(bookmark.CreatedAt.Nanosecond())},
			UpdatedAt: &timestamppb.Timestamp{Seconds: int64(bookmark.UpdatedAt.Unix()), Nanos: int32(bookmark.UpdatedAt.Nanosecond())},
		})
	}

	return &res, nil
}

func (server *bookmarkGRPCServer) ListBookmarksByArticleID(ctx context.Context, req *pb.ListBookmarksByArticleIDRequest) (*pb.ListBookmarksByArticleIDResponse, error) {
	res := pb.ListBookmarksByArticleIDResponse{}
	bookmarkRes, err := server.usecase.ListBookmarksByArticleID(ctx, int(req.ArticleId))
	if err != nil {
		return nil, err
	}
	for _, bookmark := range bookmarkRes {
		res.Bookmarks = append(res.Bookmarks, &pb.Bookmark{
			Id:        int32(bookmark.ID),
			UserId:    int32(bookmark.UserID),
			ArticleId: int32(bookmark.ArticleID),
			CreatedAt: &timestamppb.Timestamp{Seconds: int64(bookmark.CreatedAt.Unix()), Nanos: int32(bookmark.CreatedAt.Nanosecond())},
			UpdatedAt: &timestamppb.Timestamp{Seconds: int64(bookmark.UpdatedAt.Unix()), Nanos: int32(bookmark.UpdatedAt.Nanosecond())},
		})
	}

	return &res, nil
}

func (server *bookmarkGRPCServer) DeleteBookmark(ctx context.Context, req *pb.DeleteBookmarkRequest) (*pb.DeleteBookmarkResponse, error) {
	res := pb.DeleteBookmarkResponse{}

	bookmarkRes, err := server.usecase.DeleteBookmark(
		ctx,
		domain.Bookmark{
			UserID:    uint(req.UserId),
			ArticleID: uint(req.ArticleId),
		},
	)
	res.Success = bookmarkRes

	return &res, err
}

func (server *bookmarkGRPCServer) DeleteBookmarkByUserID(ctx context.Context, req *pb.DeleteBookmarkByUserIDRequest) (*pb.DeleteBookmarkByUserIDResponse, error) {
	res := pb.DeleteBookmarkByUserIDResponse{}
	bookmarkRes, err := server.usecase.DeleteBookmarkByUserID(ctx, int(req.UserId))
	res.Success = bookmarkRes

	return &res, err
}

func (server *bookmarkGRPCServer) DeleteBookmarkByUserIDCompensate(ctx context.Context, req *pb.DeleteBookmarkByUserIDRequest) (*pb.DeleteBookmarkByUserIDResponse, error) {
	res := pb.DeleteBookmarkByUserIDResponse{}
	bookmarkRes, err := server.usecase.DeleteBookmarkByUserIDCompensate(ctx, int(req.UserId))
	res.Success = bookmarkRes

	return &res, err
}

func (server *bookmarkGRPCServer) DeleteBookmarkByArticleID(ctx context.Context, req *pb.DeleteBookmarkByArticleIDRequest) (*pb.DeleteBookmarkByArticleIDResponse, error) {
	res := pb.DeleteBookmarkByArticleIDResponse{}
	bookmarkRes, err := server.usecase.DeleteBookmarkByArticleID(ctx, int(req.ArticleId))
	res.Success = bookmarkRes

	return &res, err
}

func (server *bookmarkGRPCServer) DeleteBookmarkByArticleIDCompensate(ctx context.Context, req *pb.DeleteBookmarkByArticleIDRequest) (*pb.DeleteBookmarkByArticleIDResponse, error) {
	res := pb.DeleteBookmarkByArticleIDResponse{}
	bookmarkRes, err := server.usecase.DeleteBookmarkByArticleIDCompensate(ctx, int(req.ArticleId))
	res.Success = bookmarkRes

	return &res, err
}
