//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/loak155/techbranch-backend/internal/auth/router"
	"github.com/loak155/techbranch-backend/internal/auth/usecase"
	"github.com/loak155/techbranch-backend/internal/pkg/config"
	"github.com/loak155/techbranch-backend/internal/pkg/jwt"
	"github.com/loak155/techbranch-backend/internal/user/client"
	pb "github.com/loak155/techbranch-backend/proto/auth"
	"google.golang.org/grpc"
)

func InitServer(conf *config.Config, grpcServer *grpc.Server) (pb.AuthServiceServer, error) {
	panic(wire.Build(
		client.NewUserGRPCClient,
		jwt.NewJwtManager,
		usecase.NewAuthUsecase,
		router.NewAuthGRPCServer,
	))
}
