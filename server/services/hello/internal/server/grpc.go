package server

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/justblue/luoye/services/hello/internal/conf"
	grpcserver "github.com/justblue/luoye/services/hello/internal/grpchandler"
)

func NewGRPCServer(cfg *conf.Config, greeter *grpcserver.GreeterServer) *kratosgrpc.Server {
	srv := kratosgrpc.NewServer(
		kratosgrpc.Address(cfg.GRPC.Addr),
		kratosgrpc.Middleware(
			recovery.Recovery(),
		),
	)
	greeter.Register(srv.Server)
	return srv
}
