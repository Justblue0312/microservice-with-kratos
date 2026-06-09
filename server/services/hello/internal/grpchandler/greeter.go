package grpchandler

import (
	"context"

	helloworldv1 "github.com/justblue/luoye/gen/go/helloworld"
	"github.com/justblue/luoye/services/hello/internal/domain"
	"google.golang.org/grpc"
)

type GreeterServer struct {
    helloworldv1.UnimplementedGreeterServer
    svc domain.Greeter
}

func NewGreeterServer(svc domain.Greeter) *GreeterServer {
    return &GreeterServer{svc: svc}
}

func (s *GreeterServer) SayHello(ctx context.Context, req *helloworldv1.HelloRequest) (*helloworldv1.HelloReply, error) {
    reply, err := s.svc.Greet(ctx, &domain.GreetRequest{Name: req.Name})
    if err != nil {
        return nil, err
    }
    return &helloworldv1.HelloReply{Message: reply.Message}, nil
}

func (s *GreeterServer) Register(srv *grpc.Server) {
    helloworldv1.RegisterGreeterServer(srv, s)
}
