package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-openapi/runtime/middleware"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/loak155/techbranch-backend/internal/adapter"
	"github.com/loak155/techbranch-backend/pkg/auth"
	"github.com/loak155/techbranch-backend/pkg/config"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/pkg/logger"
	"github.com/loak155/techbranch-backend/pkg/migration"
	"github.com/loak155/techbranch-backend/pkg/pb"
	"github.com/loak155/techbranch-backend/pkg/redis"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if os.Getenv("ENV") == "local" {
		if err := godotenv.Load(); err != nil {
			log.Fatal().Err(err).Msg("failed to load .env file")
		}
	}

	conf, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	waitGroup, ctx := errgroup.WithContext(ctx)

	migration.DBMigrate(conf.MigrationUrl, conf.DbSource)
	runGatewayServer(ctx, waitGroup, conf)
	runGrpcServer(ctx, waitGroup, conf)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runGrpcServer(ctx context.Context, waitGroup *errgroup.Group, conf *config.Config) {
	grpcServer, _, _, _, _, _ := adapter.NewGRPCServer(conf)

	listener, err := net.Listen("tcp", conf.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen to address")
	}

	waitGroup.Go(func() error {
		log.Info().Msg("start gRPC server")
		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Fatal().Err(err).Msg("failed to start gRPC server")
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown gRPC server")
		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server is stopped")
		return nil
	})
}

func runGatewayServer(ctx context.Context, waitGroup *errgroup.Group, conf *config.Config) {
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)

	_, articleServer, userServer, bookmarkServer, commentServer, authServer := adapter.NewGRPCServer(conf)
	if err := pb.RegisterArticleServiceHandlerServer(ctx, grpcMux, articleServer); err != nil {
		log.Fatal().Err(err).Msg("failed to register article service handler")
	}
	if err := pb.RegisterUserServiceHandlerServer(ctx, grpcMux, userServer); err != nil {
		log.Fatal().Err(err).Msg("failed to register user service handler")
	}
	if err := pb.RegisterBookmarkServiceHandlerServer(ctx, grpcMux, bookmarkServer); err != nil {
		log.Fatal().Err(err).Msg("failed to register bookmark service handler")
	}
	if err := pb.RegisterCommentServiceHandlerServer(ctx, grpcMux, commentServer); err != nil {
		log.Fatal().Err(err).Msg("failed to register comment service handler")
	}
	if err := pb.RegisterAuthServiceHandlerServer(ctx, grpcMux, authServer); err != nil {
		log.Fatal().Err(err).Msg("failed to register auth service handler")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	mux.Handle("/docs/swagger/techbranch.swagger.json", http.FileServer(http.Dir(".")))
	mux.Handle("/docs", middleware.SwaggerUI(middleware.SwaggerUIOpts{
		SpecURL: "/docs/swagger/techbranch.swagger.json",
	}, nil))

	jwtAccessTokenManager := jwt.NewJwtManager(conf.JWTIssuer, conf.JwtSecret, conf.AccessTokenExpires)
	redisAccessTokenManager := redis.NewRedisManager(conf.RedisAddress, conf.RedisAccessTokenDB, conf.AccessTokenExpires)
	authHandler := auth.NewAuthHandler(*jwtAccessTokenManager, *redisAccessTokenManager, auth.AuthRequests)

	httpServer := &http.Server{
		Addr:    conf.HttpServerAddress,
		Handler: logger.HttpLogger(authHandler.HttpAuth(mux)),
	}

	waitGroup.Go(func() error {
		log.Info().Msg("start HTTP gateway server")
		err := httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Fatal().Err(err).Msg("HTTP gateway server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown HTTP gateway server")
		err := httpServer.Shutdown(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to shutdown HTTP gateway server")
			return err
		}
		log.Info().Msg("HTTP gateway server is stopped")
		return nil
	})
}
