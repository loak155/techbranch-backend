//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/loak155/techbranch-backend/pkg/config"
	"github.com/loak155/techbranch-backend/pkg/db"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/services/comment/internal/repository"
	"github.com/loak155/techbranch-backend/services/comment/internal/router"
	"github.com/loak155/techbranch-backend/services/comment/internal/usecase"
	pb "github.com/loak155/techbranch-backend/services/comment/proto"
	"google.golang.org/grpc"
)

func InitServer(conf *config.Config, grpcServer *grpc.Server) (pb.CommentServiceServer, error) {
	panic(wire.Build(
		db.NewCommentDB,
		repository.NewCommentRepository,
		usecase.NewCommentUsecase,
		jwt.NewJwtManager,
		router.NewCommentGRPCServer,
	))
}
