package server

import (
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/justblue/luoye/services/goodbye/internal/conf"
	grpcserver "github.com/justblue/luoye/services/goodbye/internal/grpchandler"
)

func NewGRPCServer(cfg *conf.Config, goodbye *grpcserver.GoodbyeServer) *kratosgrpc.Server {
	srv := kratosgrpc.NewServer(
		kratosgrpc.Address(cfg.GRPC.Addr),
	)
	goodbye.Register(srv.Server)
	return srv
}
