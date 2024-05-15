//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/loak155/techbranch-backend/pkg/config"
	"github.com/loak155/techbranch-backend/pkg/db"
	"github.com/loak155/techbranch-backend/services/article/internal/repository"
	"github.com/loak155/techbranch-backend/services/article/internal/router"
	"github.com/loak155/techbranch-backend/services/article/internal/usecase"
	pb "github.com/loak155/techbranch-backend/services/article/proto"
	"github.com/loak155/techbranch-backend/services/bookmark/client"
	"google.golang.org/grpc"
)

func InitServer(conf *config.Config, grpcServer *grpc.Server) (pb.ArticleServiceServer, error) {
	panic(wire.Build(
		db.NewArticleDB,
		repository.NewArticleRepository,
		client.NewBookmarkGRPCClient,
		usecase.NewArticleUsecase,
		router.NewArticleGRPCServer,
	))
}
