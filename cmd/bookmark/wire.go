//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/loak155/techbranch-backend/internal/article/client"
	"github.com/loak155/techbranch-backend/internal/bookmark/repository"
	"github.com/loak155/techbranch-backend/internal/bookmark/router"
	"github.com/loak155/techbranch-backend/internal/bookmark/usecase"
	"github.com/loak155/techbranch-backend/internal/pkg/config"
	"github.com/loak155/techbranch-backend/internal/pkg/db"
	pb "github.com/loak155/techbranch-backend/proto/bookmark"
	"google.golang.org/grpc"
)

func InitServer(conf *config.Config, grpcServer *grpc.Server) (pb.BookmarkServiceServer, error) {
	panic(wire.Build(
		db.NewBookmarkDB,
		repository.NewBookmarkRepository,
		client.NewArticleGRPCClient,
		usecase.NewBookmarkUsecase,
		router.NewBookmarkGRPCServer,
	))
}
