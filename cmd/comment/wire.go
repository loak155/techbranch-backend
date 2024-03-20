//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/loak155/techbranch-backend/internal/comment/repository"
	"github.com/loak155/techbranch-backend/internal/comment/router"
	"github.com/loak155/techbranch-backend/internal/comment/usecase"
	"github.com/loak155/techbranch-backend/internal/pkg/config"
	"github.com/loak155/techbranch-backend/internal/pkg/db"
	pb "github.com/loak155/techbranch-backend/proto/comment"
	"google.golang.org/grpc"
)

func InitServer(conf *config.Config, grpcServer *grpc.Server) (pb.CommentServiceServer, error) {
	panic(wire.Build(
		db.NewCommentDB,
		repository.NewCommentRepository,
		usecase.NewCommentUsecase,
		router.NewCommentGRPCServer,
	))
}
