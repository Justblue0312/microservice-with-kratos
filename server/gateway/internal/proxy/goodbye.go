package proxy

import (
	goodbyev1 "github.com/justblue/luoye/gen/go/goodbye"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GoodbyeProxy struct {
	client goodbyev1.GoodbyeClient
}

func NewGoodbyeProxy(addr string) (*GoodbyeProxy, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return &GoodbyeProxy{client: goodbyev1.NewGoodbyeClient(conn)}, nil
}

func (p *GoodbyeProxy) Client() goodbyev1.GoodbyeClient {
	return p.client
}
