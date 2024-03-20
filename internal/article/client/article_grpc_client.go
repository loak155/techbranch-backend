package client

import (
	"fmt"

	"github.com/loak155/techbranch-backend/internal/pkg/config"
	pb "github.com/loak155/techbranch-backend/proto/article"
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

// type IArticleGRPCClient interface {
// 	GetArticle(req *pb.GetArticleRequest) (*pb.GetArticleResponse, error)
// 	UpdateArticle(req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error)
// }

// type articleGRPCClient struct {
// 	client pb.ArticleServiceClient
// }

// func NewArticleGRPCClient() (IArticleGRPCClient, error) {
// 	address := fmt.Sprintf("%s:%s", os.Getenv("ARTICLE_SERVICE_HOST"), os.Getenv("ARTICLE_SERVICE_PORT"))
// 	conn, err := grpc.Dial(
// 		address,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithBlock(),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	client := pb.NewArticleServiceClient(conn)

// 	return &articleGRPCClient{client}, nil
// }

// func (c *articleGRPCClient) GetArticle(req *pb.GetArticleRequest) (*pb.GetArticleResponse, error) {
// 	res, err := c.client.GetArticle(context.Background(), req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return res, nil
// }

// func (c *articleGRPCClient) UpdateArticle(req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error) {
// 	res, err := c.client.UpdateArticle(context.Background(), req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res, nil
// }
