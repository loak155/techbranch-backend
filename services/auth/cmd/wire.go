//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/loak155/techbranch-backend/pkg/config"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/pkg/oauth"
	"github.com/loak155/techbranch-backend/services/auth/internal/router"
	"github.com/loak155/techbranch-backend/services/auth/internal/usecase"
	pb "github.com/loak155/techbranch-backend/services/auth/proto"
	"github.com/loak155/techbranch-backend/services/user/client"
	"google.golang.org/grpc"
)

func InitServer(conf *config.Config, grpcServer *grpc.Server) (pb.AuthServiceServer, error) {
	panic(wire.Build(
		client.NewUserGRPCClient,
		oauth.NewGoogle,
		jwt.NewJwtManager,
		usecase.NewAuthUsecase,
		router.NewAuthGRPCServer,
	))
}
