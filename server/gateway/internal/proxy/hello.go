package proxy

import (
	helloworldv1 "github.com/justblue/luoye/gen/go/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HelloProxy struct {
	client helloworldv1.GreeterClient
}

func NewHelloProxy(addr string) (*HelloProxy, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return &HelloProxy{client: helloworldv1.NewGreeterClient(conn)}, nil
}

func (p *HelloProxy) Client() helloworldv1.GreeterClient {
	return p.client
}
