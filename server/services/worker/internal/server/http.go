package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/justblue/luoye/services/worker/internal/conf"
)

func NewHTTPServer(cfg *conf.Config) *kratoshttp.Server {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	srv := kratoshttp.NewServer(
		kratoshttp.Address(cfg.HTTP.Addr),
	)
	srv.HandlePrefix("/", r)
	return srv
}
