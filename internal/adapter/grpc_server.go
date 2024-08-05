package adapter

import (
	"github.com/loak155/techbranch-backend/internal/repository"
	"github.com/loak155/techbranch-backend/internal/usecase"
	"github.com/loak155/techbranch-backend/pkg/config"
	"github.com/loak155/techbranch-backend/pkg/db"
	"github.com/loak155/techbranch-backend/pkg/logger"
	"github.com/loak155/techbranch-backend/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(conf *config.Config) (*grpc.Server, pb.ArticleServiceServer, pb.UserServiceServer, pb.BookmarkServiceServer, pb.CommentServiceServer) {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.GrpcLogger),
	)

	gormDB := db.NewDB(conf.DbSource)

	articleRepository := repository.NewArticleRepository(gormDB)
	articleUsecase := usecase.NewArticleUsecase(articleRepository)
	articleServer := NewArticleGRPCServer(grpcServer, articleUsecase)

	userRepository := repository.NewUserRepository(gormDB)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userServer := NewUserGRPCServer(grpcServer, userUsecase)

	bookmarkRepository := repository.NewBookmarkRepository(gormDB)
	bookmarkUsecase := usecase.NewBookmarkUsecase(bookmarkRepository)
	bookmarkServer := NewBookmarkGRPCServer(grpcServer, bookmarkUsecase)

	commentRepository := repository.NewCommentRepository(gormDB)
	commentUsecase := usecase.NewCommentUsecase(commentRepository)
	commentServer := NewCommentGRPCServer(grpcServer, commentUsecase)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("grpc-server", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)
	return grpcServer, articleServer, userServer, bookmarkServer, commentServer
}
