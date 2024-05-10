//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/loak155/techbranch-backend/pkg/config"
	"github.com/loak155/techbranch-backend/pkg/db"
	"github.com/loak155/techbranch-backend/services/user/internal/repository"
	"github.com/loak155/techbranch-backend/services/user/internal/router"
	"github.com/loak155/techbranch-backend/services/user/internal/usecase"
	pb "github.com/loak155/techbranch-backend/services/user/proto"
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
