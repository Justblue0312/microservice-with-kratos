package main

import (
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/justblue/luoye/services/hello/internal/conf"
	"github.com/justblue/luoye/services/hello/internal/event"
	"github.com/justblue/luoye/services/hello/internal/grpcclient"
	grpcserver "github.com/justblue/luoye/services/hello/internal/grpchandler"
	"github.com/justblue/luoye/services/hello/internal/server"
	"github.com/justblue/luoye/services/hello/internal/usecase"
)

func initApp(cfg *conf.Config) (*kratos.App, error) {
	goodbyeClient, err := grpcclient.NewGoodbyeClient(cfg.Upstream.Goodbye)
	if err != nil {
		return nil, err
	}
	svc := usecase.NewGreeterService(goodbyeClient)
	grpcServer := server.NewGRPCServer(cfg, grpcserver.NewGreeterServer(svc))

	consumer, err := event.NewConsumer(cfg.NATS.URL)
	if err != nil {
		return nil, err
	}

	return server.NewApp(grpcServer, consumer), nil
}

func main() {
	cfg := &conf.Config{
		GRPC:     conf.GRPCConfig{Addr: ":9081"},
		Upstream: conf.UpstreamConfig{Goodbye: "localhost:9082"},
		NATS:     conf.NATSConfig{URL: "nats://localhost:4222"},
	}
	app, err := initApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
