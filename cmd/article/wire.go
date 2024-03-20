//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/loak155/techbranch-backend/internal/article/repository"
	"github.com/loak155/techbranch-backend/internal/article/router"
	"github.com/loak155/techbranch-backend/internal/article/usecase"
	"github.com/loak155/techbranch-backend/internal/pkg/config"
	"github.com/loak155/techbranch-backend/internal/pkg/db"
	pb "github.com/loak155/techbranch-backend/proto/article"
	"google.golang.org/grpc"
)

func InitServer(conf *config.Config, grpcServer *grpc.Server) (pb.ArticleServiceServer, error) {
	panic(wire.Build(
		db.NewArticleDB,
		repository.NewArticleRepository,
		usecase.NewArticleUsecase,
		router.NewArticleGRPCServer,
	))
}
