package server

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	helloworldv1 "github.com/justblue/luoye/gen/go/helloworld"
	"github.com/justblue/luoye/services/hello/internal/conf"
	grpchandler "github.com/justblue/luoye/services/hello/internal/grpchandler"
)

func NewHTTPServer(cfg *conf.Config, greeter *grpchandler.GreeterServer) *kratoshttp.Server {
	srv := kratoshttp.NewServer(
		kratoshttp.Address(cfg.HTTP.Addr),
		kratoshttp.Middleware(
			recovery.Recovery(),
		),
	)
	helloworldv1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
