package server

import (
	"encoding/json"
	"net/http"

	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	helloworldv1 "github.com/justblue/luoye/gen/go/helloworld"
	"github.com/justblue/luoye/gateway/internal/conf"
	"github.com/justblue/luoye/gateway/internal/proxy"
)

func NewHTTPServer(cfg *conf.Config, hello *proxy.HelloProxy) *kratoshttp.Server {
	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Get("/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		reply, err := hello.Client().SayHello(r.Context(), &helloworldv1.HelloRequest{Name: name})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reply)
	})

	srv := kratoshttp.NewServer(kratoshttp.Address(cfg.HTTP.Addr))
	srv.HandlePrefix("/", r)
	return srv
}
