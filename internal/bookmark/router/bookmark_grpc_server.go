package router

import (
	"context"

	"github.com/loak155/techbranch-backend/internal/bookmark/domain"
	"github.com/loak155/techbranch-backend/internal/bookmark/usecase"
	pb "github.com/loak155/techbranch-backend/proto/bookmark"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IBookmarkGRPCServer interface {
	CreateBookmark(ctx context.Context, req *pb.CreateBookmarkRequest) (*pb.CreateBookmarkResponse, error)
	GetBookmark(context.Context, *pb.GetBookmarkRequest) (*pb.GetBookmarkResponse, error)
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
	bu usecase.IBookmarkUsecase
}

func NewBookmarkGRPCServer(grpcServer *grpc.Server, bu usecase.IBookmarkUsecase) pb.BookmarkServiceServer {
	s := bookmarkGRPCServer{bu: bu}
	pb.RegisterBookmarkServiceServer(grpcServer, &s)
	reflection.Register(grpcServer)
	return &s
}

func (s *bookmarkGRPCServer) CreateBookmark(ctx context.Context, req *pb.CreateBookmarkRequest) (*pb.CreateBookmarkResponse, error) {
	res := pb.CreateBookmarkResponse{}
	bookmarkRes, err := s.bu.CreateBookmark(
		ctx,
		domain.Bookmark{
			UserID:    uint(req.Bookmark.UserId),
			ArticleID: uint(req.Bookmark.ArticleId),
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

func (s *bookmarkGRPCServer) GetBookmark(ctx context.Context, req *pb.GetBookmarkRequest) (*pb.GetBookmarkResponse, error) {
	res := pb.GetBookmarkResponse{}
	bookmarkRes, err := s.bu.GetBookmark(ctx, int(req.Id))
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

func (s *bookmarkGRPCServer) ListBookmarks(ctx context.Context, req *pb.ListBookmarksRequest) (*pb.ListBookmarksResponse, error) {
	res := pb.ListBookmarksResponse{}
	bookmarkRes, err := s.bu.ListBookmarks(ctx)
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

func (s *bookmarkGRPCServer) ListBookmarksByUserID(ctx context.Context, req *pb.ListBookmarksByUserIDRequest) (*pb.ListBookmarksByUserIDResponse, error) {
	res := pb.ListBookmarksByUserIDResponse{}
	bookmarkRes, err := s.bu.ListBookmarksByUserID(ctx, int(req.UserId))
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

func (s *bookmarkGRPCServer) ListBookmarksByArticleID(ctx context.Context, req *pb.ListBookmarksByArticleIDRequest) (*pb.ListBookmarksByArticleIDResponse, error) {
	res := pb.ListBookmarksByArticleIDResponse{}
	bookmarkRes, err := s.bu.ListBookmarksByArticleID(ctx, int(req.ArticleId))
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

func (s *bookmarkGRPCServer) DeleteBookmark(ctx context.Context, req *pb.DeleteBookmarkRequest) (*pb.DeleteBookmarkResponse, error) {
	res := pb.DeleteBookmarkResponse{}
	bookmarkRes, err := s.bu.DeleteBookmark(
		ctx,
		domain.Bookmark{
			UserID:    uint(req.Bookmark.UserId),
			ArticleID: uint(req.Bookmark.ArticleId),
		},
	)
	res.Success = bookmarkRes

	return &res, err
}

func (s *bookmarkGRPCServer) DeleteBookmarkByUserID(ctx context.Context, req *pb.DeleteBookmarkByUserIDRequest) (*pb.DeleteBookmarkByUserIDResponse, error) {
	res := pb.DeleteBookmarkByUserIDResponse{}
	bookmarkRes, err := s.bu.DeleteBookmarkByUserID(ctx, int(req.UserId))
	res.Success = bookmarkRes

	return &res, err
}

func (s *bookmarkGRPCServer) DeleteBookmarkByUserIDCompensate(ctx context.Context, req *pb.DeleteBookmarkByUserIDRequest) (*pb.DeleteBookmarkByUserIDResponse, error) {
	res := pb.DeleteBookmarkByUserIDResponse{}
	bookmarkRes, err := s.bu.DeleteBookmarkByUserIDCompensate(ctx, int(req.UserId))
	res.Success = bookmarkRes

	return &res, err
}

func (s *bookmarkGRPCServer) DeleteBookmarkByArticleID(ctx context.Context, req *pb.DeleteBookmarkByArticleIDRequest) (*pb.DeleteBookmarkByArticleIDResponse, error) {
	res := pb.DeleteBookmarkByArticleIDResponse{}
	bookmarkRes, err := s.bu.DeleteBookmarkByArticleID(ctx, int(req.ArticleId))
	res.Success = bookmarkRes

	return &res, err
}

func (s *bookmarkGRPCServer) DeleteBookmarkByArticleIDCompensate(ctx context.Context, req *pb.DeleteBookmarkByArticleIDRequest) (*pb.DeleteBookmarkByArticleIDResponse, error) {
	res := pb.DeleteBookmarkByArticleIDResponse{}
	bookmarkRes, err := s.bu.DeleteBookmarkByArticleIDCompensate(ctx, int(req.ArticleId))
	res.Success = bookmarkRes

	return &res, err
}
