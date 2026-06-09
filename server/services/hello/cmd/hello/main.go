package main

import (
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/justblue/luoye/services/hello/internal/conf"
	grpcserver "github.com/justblue/luoye/services/hello/internal/grpchandler"
	"github.com/justblue/luoye/services/hello/internal/httphandler"
	"github.com/justblue/luoye/services/hello/internal/server"
	"github.com/justblue/luoye/services/hello/internal/usecase"
)

func initApp(cfg *conf.Config) (*kratos.App, error) {
	svc := usecase.NewGreeterService()
	httpServer := server.NewHTTPServer(cfg, httphandler.NewGreeterHandler(svc))
	grpcServer := server.NewGRPCServer(cfg, grpcserver.NewGreeterServer(svc))
	return server.NewApp(httpServer, grpcServer), nil
}

func main() {
	cfg := &conf.Config{
		HTTP: conf.HTTPConfig{Addr: ":8081"},
		GRPC: conf.GRPCConfig{Addr: ":9081"},
	}
	app, err := initApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
