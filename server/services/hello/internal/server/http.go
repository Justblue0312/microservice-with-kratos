package server

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/justblue/luoye/services/hello/internal/conf"
	"github.com/justblue/luoye/services/hello/internal/httphandler"
)

func NewHTTPServer(cfg *conf.Config, greeter *httphandler.GreeterHandler) *kratoshttp.Server {
    // Build a Chi router
    r := chi.NewRouter()
    r.Use(chiMiddleware.RequestID)
    r.Use(chiMiddleware.Logger)
    r.Use(chiMiddleware.Recoverer)

    r.Route("/v1", greeter.Routes())

    // Wrap Chi inside a Kratos HTTP server for lifecycle management
    srv := kratoshttp.NewServer(
        kratoshttp.Address(cfg.HTTP.Addr),
    )
    srv.HandlePrefix("/", r)
    return srv
}
