package main

import (
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/justblue/luoye/services/goodbye/internal/conf"
	"github.com/justblue/luoye/services/goodbye/internal/event"
	grpcserver "github.com/justblue/luoye/services/goodbye/internal/grpchandler"
	"github.com/justblue/luoye/services/goodbye/internal/server"
	"github.com/justblue/luoye/services/goodbye/internal/usecase"
)

func initApp(cfg *conf.Config) (*kratos.App, error) {
	publisher, err := event.NewPublisher(cfg.NATS.URL)
	if err != nil {
		return nil, err
	}
	svc := usecase.NewGoodbyeService(publisher)
	grpcServer := server.NewGRPCServer(cfg, grpcserver.NewGoodbyeServer(svc))
	return server.NewApp(grpcServer), nil
}

func main() {
	cfg := &conf.Config{
		GRPC: conf.GRPCConfig{Addr: ":9082"},
		NATS: conf.NATSConfig{URL: "nats://localhost:4222"},
	}
	app, err := initApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
