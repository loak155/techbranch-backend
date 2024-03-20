//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/loak155/techbranch-backend/internal/pkg/config"
	"github.com/loak155/techbranch-backend/internal/pkg/db"
	"github.com/loak155/techbranch-backend/internal/user/repository"
	"github.com/loak155/techbranch-backend/internal/user/router"
	"github.com/loak155/techbranch-backend/internal/user/usecase"
	pb "github.com/loak155/techbranch-backend/proto/user"
	"google.golang.org/grpc"
)

func InitServer(conf *config.Config, grpcServer *grpc.Server) (pb.UserServiceServer, error) {
	panic(wire.Build(
		db.NewUserDB,
		repository.NewUserRepository,
		usecase.NewUserUsecase,
		router.NewUserGRPCServer,
	))
}
