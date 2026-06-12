package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/justblue/luoye/gateway/internal/conf"
	"github.com/justblue/luoye/gateway/internal/proxy"
	goodbyev1 "github.com/justblue/luoye/gen/go/goodbye"
	helloworldv1 "github.com/justblue/luoye/gen/go/helloworld"
)

func NewHTTPServer(cfg *conf.Config, hello *proxy.HelloProxy, goodbye *proxy.GoodbyeProxy) *kratoshttp.Server {
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

	r.Get("/v1/goodbye", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		reply, err := goodbye.Client().SayGoodbye(r.Context(), &goodbyev1.GoodbyeRequest{Name: name})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reply)
	})

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
	return srv
}
