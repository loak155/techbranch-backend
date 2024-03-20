package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/loak155/techbranch-backend/internal/pkg/config"
	"github.com/loak155/techbranch-backend/internal/pkg/interceptor"
	"github.com/loak155/techbranch-backend/internal/pkg/jwt"
	"google.golang.org/grpc"
)

var flagConfig = flag.String("config", "./configs/config.yaml", "path to config file")

func main() {
	slog.Info("starting comment service")

	flag.Parse()
	conf, err := config.Load(*flagConfig)
	if err != nil {
		slog.Error("failed to load config: ", err)
	}

	authInterceptor := interceptor.NewAuthInterceptor(jwt.NewJwtManager(conf))

	ctx, cancel := context.WithCancel(context.Background())
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggingInterceptor,
			auth.UnaryServerInterceptor(authInterceptor.AuthFunc),
		),
	)
	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()

	InitServer(conf, server)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Comment.Server.Port))
	if err != nil {
		slog.Error("failed to listen to address")
		cancel()
	}
	err = server.Serve(listener)
	if err != nil {
		slog.Error("failed to start gRPC server")
		cancel()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case v := <-quit:
		slog.Info("signal.Notify: ", v)
	case done := <-ctx.Done():
		slog.Info("ctx.Done: ", done)
	}
}
