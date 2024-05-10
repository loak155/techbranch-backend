package interceptor

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	slog.Info("[Message]", "Method", info.FullMethod, "Request", req)
	res, err := handler(ctx, req)
	slog.Info("[Message]", "Response", res, "Error", err)
	return res, err
}
