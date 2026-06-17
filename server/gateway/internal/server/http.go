package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/justblue/luoye/gateway/internal/conf"
	"github.com/justblue/luoye/gateway/internal/proxy"
)

func NewHTTPServer(cfg *conf.Config) (*kratoshttp.Server, error) {
	helloProxy, err := proxy.NewReverseProxy(cfg.Upstreams.Hello)
	if err != nil {
		return nil, err
	}
	goodbyeProxy, err := proxy.NewReverseProxy(cfg.Upstreams.Goodbye)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Mount("/v1/hello", http.StripPrefix("/v1/hello", helloProxy))
	r.Mount("/v1/goodbye", http.StripPrefix("/v1/goodbye", goodbyeProxy))

	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		openapiPath := filepath.Join("gen", "openapi", "openapi.yaml")
		data, err := os.ReadFile(openapiPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Write(data)
	})

	srv := kratoshttp.NewServer(kratoshttp.Address(cfg.HTTP.Addr))
	srv.HandlePrefix("/", r)
	return srv, nil
}
