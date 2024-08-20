package adapter

import (
	"github.com/loak155/techbranch-backend/internal/repository"
	"github.com/loak155/techbranch-backend/internal/usecase"
	"github.com/loak155/techbranch-backend/pkg/auth"
	"github.com/loak155/techbranch-backend/pkg/config"
	"github.com/loak155/techbranch-backend/pkg/db"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/pkg/logger"
	"github.com/loak155/techbranch-backend/pkg/mail"
	"github.com/loak155/techbranch-backend/pkg/oauth"
	"github.com/loak155/techbranch-backend/pkg/pb"
	"github.com/loak155/techbranch-backend/pkg/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(conf *config.Config) (*grpc.Server, pb.ArticleServiceServer, pb.UserServiceServer, pb.BookmarkServiceServer, pb.CommentServiceServer, pb.AuthServiceServer) {
	jwtAccessTokenManager := jwt.NewJwtManager(conf.JWTIssuer, conf.JwtSecret, conf.AccessTokenExpires)
	jwtRefreshTokenManager := jwt.NewJwtManager(conf.JWTIssuer, conf.JwtSecret, conf.RefreshTokenExpires)
	redisAccessTokenManager := redis.NewRedisManager(conf.RedisAddress, conf.RedisAccessTokenDB, conf.AccessTokenExpires)
	redisRefreshTokenManager := redis.NewRedisManager(conf.RedisAddress, conf.RedisRefreshTokenDB, conf.RefreshTokenExpires)
	google := oauth.NewGoogleManager(conf.OauthGoogleState, conf.OauthGoogleClientID, conf.OauthGoogleClientSecret, conf.OauthGoogleRedirectURL)

	authInterceptor := auth.NewAuthInterceptor(*jwtAccessTokenManager, *redisAccessTokenManager, auth.AuthMethods)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logger.GrpcLogger,
			authInterceptor.Unary(),
		),
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

	presignupMailManager, _ := mail.NewPresignupMailManager(conf.GmailFrom, conf.GmailPassword, conf.PresignupMailSubject, conf.PresignupMailTemplate, conf.SignupURL)
	presignupRedisManager := redis.NewRedisManager(conf.RedisAddress, conf.RedisPresignupDB, conf.PresignupExpires)
	authUsecase := usecase.NewAuthUsecase(userRepository, *jwtAccessTokenManager, *jwtRefreshTokenManager, *redisAccessTokenManager, *redisRefreshTokenManager, *google, *presignupRedisManager, *presignupMailManager)
	authServer := NewAuthGRPCServer(grpcServer, authUsecase)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("grpc-server", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)
	return grpcServer, articleServer, userServer, bookmarkServer, commentServer, authServer
}
