package server

import (
	"net/http"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	goodbyev1 "github.com/justblue/luoye/gen/go/goodbye"
	"github.com/justblue/luoye/services/goodbye/internal/conf"
	grpchandler "github.com/justblue/luoye/services/goodbye/internal/grpchandler"
)

func NewHTTPServer(cfg *conf.Config, goodbye *grpchandler.GoodbyeServer) *kratoshttp.Server {
	srv := kratoshttp.NewServer(
		kratoshttp.Address(cfg.HTTP.Addr),
		kratoshttp.Middleware(
			recovery.Recovery(),
		),
	)
	r := srv.Route("/")
	r.GET("/health", func(ctx kratoshttp.Context) error {
		return ctx.Result(http.StatusOK, map[string]string{"status": "ok"})
	})
	goodbyev1.RegisterGoodbyeHTTPServer(srv, goodbye)
	return srv
}
