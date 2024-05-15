//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/loak155/techbranch-backend/pkg/config"
	"github.com/loak155/techbranch-backend/pkg/db"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/services/bookmark/internal/repository"
	"github.com/loak155/techbranch-backend/services/bookmark/internal/router"
	"github.com/loak155/techbranch-backend/services/bookmark/internal/usecase"
	pb "github.com/loak155/techbranch-backend/services/bookmark/proto"
	"google.golang.org/grpc"
)

func InitServer(conf *config.Config, grpcServer *grpc.Server) (pb.BookmarkServiceServer, error) {
	panic(wire.Build(
		db.NewBookmarkDB,
		repository.NewBookmarkRepository,
		usecase.NewBookmarkUsecase,
		jwt.NewJwtManager,
		router.NewBookmarkGRPCServer,
	))
}
