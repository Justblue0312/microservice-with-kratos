package grpchandler

import (
	"context"

	goodbyev1 "github.com/justblue/luoye/gen/go/goodbye"
	"github.com/justblue/luoye/services/goodbye/internal/domain"
	"google.golang.org/grpc"
)

type GoodbyeServer struct {
	goodbyev1.UnimplementedGoodbyeServer
	svc domain.Goodbye
}

func NewGoodbyeServer(svc domain.Goodbye) *GoodbyeServer {
	return &GoodbyeServer{svc: svc}
}

func (s *GoodbyeServer) SayGoodbye(ctx context.Context, req *goodbyev1.GoodbyeRequest) (*goodbyev1.GoodbyeReply, error) {
	reply, err := s.svc.SayGoodbye(ctx, &domain.GoodbyeRequest{Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &goodbyev1.GoodbyeReply{Message: reply.Message}, nil
}

func (s *GoodbyeServer) Register(srv *grpc.Server) {
	goodbyev1.RegisterGoodbyeServer(srv, s)
}
