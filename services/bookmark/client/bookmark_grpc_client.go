package client

import (
	"fmt"

	"github.com/loak155/techbranch-backend/pkg/config"
	pb "github.com/loak155/techbranch-backend/services/bookmark/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewBookmarkGRPCClient(conf *config.Config) (pb.BookmarkServiceClient, error) {
	address := fmt.Sprintf("%s:%d", conf.Bookmark.Server.Host, conf.Bookmark.Server.Port)
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}
	return pb.NewBookmarkServiceClient(conn), nil
}
