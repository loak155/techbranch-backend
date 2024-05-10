package client

import (
	"fmt"

	"github.com/loak155/techbranch-backend/pkg/config"
	pb "github.com/loak155/techbranch-backend/services/user/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserGRPCClient(conf *config.Config) (pb.UserServiceClient, error) {
	address := fmt.Sprintf("%s:%d", conf.User.Server.Host, conf.User.Server.Port)
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}
	return pb.NewUserServiceClient(conn), nil
}
