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

type IBookmarkGRPCServer interface {
	CreateBookmark(ctx context.Context, req *pb.CreateBookmarkRequest) (*pb.CreateBookmarkResponse, error)
	GetBookmarkCountByArticleID(ctx context.Context, req *pb.GetBookmarkCountByArticleIDRequest) (*pb.GetBookmarkCountByArticleIDResponse, error)
	ListBookmarksByUserID(ctx context.Context, req *pb.ListBookmarksByUserIDRequest) (*pb.ListBookmarksByUserIDResponse, error)
	ListBookmarksByArticleID(ctx context.Context, req *pb.ListBookmarksByArticleIDRequest) (*pb.ListBookmarksByArticleIDResponse, error)
	DeleteBookmarkByUserIDAndArticleID(ctx context.Context, req *pb.DeleteBookmarkByUserIDAndArticleIDRequest) (*pb.DeleteBookmarkByUserIDAndArticleIDResponse, error)
	DeleteBookmarkByUserID(ctx context.Context, req *pb.DeleteBookmarkByUserIDRequest) (*pb.DeleteBookmarkByUserIDResponse, error)
	DeleteBookmarkByArticleID(ctx context.Context, req *pb.DeleteBookmarkByArticleIDRequest) (*pb.DeleteBookmarkByArticleIDResponse, error)
}

type bookmarkGRPCServer struct {
	pb.UnimplementedBookmarkServiceServer
	usecase usecase.IBookmarkUsecase
}

func NewBookmarkGRPCServer(grpcServer *grpc.Server, usecase usecase.IBookmarkUsecase) pb.BookmarkServiceServer {
	server := bookmarkGRPCServer{usecase: usecase}
	pb.RegisterBookmarkServiceServer(grpcServer, &server)
	return &server
}

func (server *bookmarkGRPCServer) CreateBookmark(ctx context.Context, req *pb.CreateBookmarkRequest) (*pb.CreateBookmarkResponse, error) {
	res := pb.CreateBookmarkResponse{}
	bookmark, err := server.usecase.CreateBookmark(
		domain.Bookmark{
			UserID:    uint(req.UserId),
			ArticleID: uint(req.ArticleId),
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create bookmark: %v", err)
	}

	res.Bookmark = &pb.Bookmark{
		Id:        int32(bookmark.ID),
		UserId:    int32(bookmark.UserID),
		ArticleId: int32(bookmark.ArticleID),
		CreatedAt: &timestamppb.Timestamp{Seconds: int64(bookmark.CreatedAt.Unix()), Nanos: int32(bookmark.CreatedAt.Nanosecond())},
		UpdatedAt: &timestamppb.Timestamp{Seconds: int64(bookmark.UpdatedAt.Unix()), Nanos: int32(bookmark.UpdatedAt.Nanosecond())},
	}

	return &res, nil
}

func (server *bookmarkGRPCServer) GetBookmarkCountByArticleID(ctx context.Context, req *pb.GetBookmarkCountByArticleIDRequest) (*pb.GetBookmarkCountByArticleIDResponse, error) {
	res := pb.GetBookmarkCountByArticleIDResponse{}
	count, err := server.usecase.GetBookmarkCountByArticleID(int(req.ArticleId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get bookmark count by article id: %v", err)
	}
	res.Count = int32(count)
	return &res, nil
}

func (server *bookmarkGRPCServer) ListBookmarksByUserID(ctx context.Context, req *pb.ListBookmarksByUserIDRequest) (*pb.ListBookmarksByUserIDResponse, error) {
	res := pb.ListBookmarksByUserIDResponse{}
	bookmarkRes, err := server.usecase.ListBookmarksByUserID(int(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list bookmarks by user id: %v", err)
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
	// if len(bookmarkRes) == 0 {
	// 	res.Bookmarks = append(res.Bookmarks, &pb.Bookmark{})
	// }

	return &res, nil
}

func (server *bookmarkGRPCServer) ListBookmarksByArticleID(ctx context.Context, req *pb.ListBookmarksByArticleIDRequest) (*pb.ListBookmarksByArticleIDResponse, error) {
	res := pb.ListBookmarksByArticleIDResponse{}
	bookmarkRes, err := server.usecase.ListBookmarksByArticleID(int(req.ArticleId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list bookmarks by article id: %v", err)
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
	// if len(bookmarkRes) == 0 {
	// 	res.Bookmarks = append(res.Bookmarks, &pb.Bookmark{})
	// }

	return &res, nil
}

func (server *bookmarkGRPCServer) DeleteBookmarkByUserIDAndArticleID(ctx context.Context, req *pb.DeleteBookmarkByUserIDAndArticleIDRequest) (*pb.DeleteBookmarkByUserIDAndArticleIDResponse, error) {
	res := pb.DeleteBookmarkByUserIDAndArticleIDResponse{}
	err := server.usecase.DeleteBookmarkByUserIDAndArticleID(int(req.UserId), int(req.ArticleId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete bookmark by user id and article id: %v", err)
	}

	return &res, err
}

func (server *bookmarkGRPCServer) DeleteBookmarkByUserID(ctx context.Context, req *pb.DeleteBookmarkByUserIDRequest) (*pb.DeleteBookmarkByUserIDResponse, error) {
	res := pb.DeleteBookmarkByUserIDResponse{}
	err := server.usecase.DeleteBookmarkByUserID(int(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete bookmark by user id and article id: %v", err)
	}

	return &res, err
}

func (server *bookmarkGRPCServer) DeleteBookmarkByArticleID(ctx context.Context, req *pb.DeleteBookmarkByArticleIDRequest) (*pb.DeleteBookmarkByArticleIDResponse, error) {
	res := pb.DeleteBookmarkByArticleIDResponse{}
	err := server.usecase.DeleteBookmarkByArticleID(int(req.ArticleId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete bookmark by user id and article id: %v", err)
	}

	return &res, err
}
