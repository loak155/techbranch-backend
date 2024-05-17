package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/loak155/techbranch-backend/pkg/config"
	"github.com/loak155/techbranch-backend/pkg/interceptor"
	"google.golang.org/grpc"
)

func main() {
	slog.Info("starting user service")

	if os.Getenv("ENV") == "local" {
		if err := godotenv.Load(); err != nil {
			slog.Error("failed to load .env file: ", err)
		}
	}

	conf, err := config.Load()
	if err != nil {
		slog.Error("failed to load config: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggingInterceptor),
	)
	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()

	InitServer(conf, server)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.User.Server.Port))
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
