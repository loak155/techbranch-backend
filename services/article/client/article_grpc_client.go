package client

import (
	"fmt"

	"github.com/loak155/techbranch-backend/pkg/config"
	pb "github.com/loak155/techbranch-backend/services/article/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewArticleGRPCClient(conf *config.Config) (pb.ArticleServiceClient, error) {
	address := fmt.Sprintf("%s:%d", conf.Article.Server.Host, conf.Article.Server.Port)
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}
	return pb.NewArticleServiceClient(conn), nil
}
