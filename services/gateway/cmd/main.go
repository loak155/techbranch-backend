package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/loak155/techbranch-backend/pkg/config"
	articlepb "github.com/loak155/techbranch-backend/services/article/proto"
	authpb "github.com/loak155/techbranch-backend/services/auth/proto"
	bookmarkpb "github.com/loak155/techbranch-backend/services/bookmark/proto"
	commentpb "github.com/loak155/techbranch-backend/services/comment/proto"
)

var flagConfig = flag.String("config", "./configs/config.yaml", "path to config file")

func withLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Run request", "http_method", r.Method, "http_url", r.URL)

		body, _ := ioutil.ReadAll(r.Body)
		slog.Info("Request body", "body", string(body))
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		h.ServeHTTP(w, r)
	})
}

func enableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func newGateway(ctx context.Context, conf *config.Config) (http.Handler, error) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	authEndpoint := fmt.Sprintf("%s:%d", conf.Auth.Server.Host, conf.Auth.Server.Port)
	err := authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, authEndpoint, opts)
	if err != nil {
		return nil, err
	}
	articleEndpoint := fmt.Sprintf("%s:%d", conf.Article.Server.Host, conf.Article.Server.Port)
	err = articlepb.RegisterArticleServiceHandlerFromEndpoint(ctx, mux, articleEndpoint, opts)
	if err != nil {
		return nil, err
	}
	bookmarkEndpoint := fmt.Sprintf("%s:%d", conf.Bookmark.Server.Host, conf.Bookmark.Server.Port)
	err = bookmarkpb.RegisterBookmarkServiceHandlerFromEndpoint(ctx, mux, bookmarkEndpoint, opts)
	if err != nil {
		return nil, err
	}
	commentEndpoint := fmt.Sprintf("%s:%d", conf.Comment.Server.Host, conf.Comment.Server.Port)
	err = commentpb.RegisterCommentServiceHandlerFromEndpoint(ctx, mux, commentEndpoint, opts)
	if err != nil {
		return nil, err
	}

	return mux, err
}

func main() {
	slog.Info("starting gateway")

	flag.Parse()
	conf, err := config.Load(*flagConfig)
	if err != nil {
		slog.Error("failed to load config: ", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux, err := newGateway(ctx, conf)
	if err != nil {
		slog.Error("failed to create a new gateway", err)
	}

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.Gateway.Server.Host, conf.Gateway.Server.Port),
		Handler: withLogger(enableCors(mux)),
	}
	go func() {
		defer s.Close()
		<-ctx.Done()
	}()

	s.ListenAndServe()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case v := <-quit:
		slog.Info("signal.Notify: ", v)
	case done := <-ctx.Done():
		slog.Info("ctx.Done: ", done)
	}
}
