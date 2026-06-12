package grpcclient

import (
	"context"

	goodbyev1 "github.com/justblue/luoye/gen/go/goodbye"
	"github.com/justblue/luoye/services/hello/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GoodbyeClient struct {
	client goodbyev1.GoodbyeClient
}

func NewGoodbyeClient(addr string) (*GoodbyeClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return &GoodbyeClient{client: goodbyev1.NewGoodbyeClient(conn)}, nil
}

func (c *GoodbyeClient) SayGoodbye(ctx context.Context, req *domain.GoodbyeRequest) (*domain.GoodbyeReply, error) {
	reply, err := c.client.SayGoodbye(ctx, &goodbyev1.GoodbyeRequest{Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &domain.GoodbyeReply{Message: reply.Message}, nil
}
