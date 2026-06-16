package grpcclient

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	goodbyev1 "github.com/justblue/luoye/gen/go/goodbye"
	"github.com/justblue/luoye/services/hello/internal/conf"
	"github.com/justblue/luoye/services/hello/internal/domain"
	"google.golang.org/grpc"
)

type GoodbyeClient struct {
	client goodbyev1.GoodbyeClient
	conn   *grpc.ClientConn
}

func NewGoodbyeClient(upstream conf.UpstreamConfig) (*GoodbyeClient, error) {
	conn, err := kratosgrpc.DialInsecure(
		context.Background(),
		kratosgrpc.WithEndpoint(upstream.Goodbye),
		kratosgrpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		return nil, err
	}
	return &GoodbyeClient{client: goodbyev1.NewGoodbyeClient(conn), conn: conn}, nil
}

func (c *GoodbyeClient) SayGoodbye(ctx context.Context, req *domain.GoodbyeRequest) (*domain.GoodbyeReply, error) {
	reply, err := c.client.SayGoodbye(ctx, &goodbyev1.GoodbyeRequest{Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &domain.GoodbyeReply{Message: reply.Message}, nil
}

func (c *GoodbyeClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
