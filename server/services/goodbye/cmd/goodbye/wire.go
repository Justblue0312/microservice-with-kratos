//go:build wireinject

package main

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	"github.com/hibiken/asynq"
	kitasynq "github.com/justblue/luoye/kit/messaging/asynq"
	kitnats "github.com/justblue/luoye/kit/messaging/nats"
	"github.com/justblue/luoye/services/goodbye/internal/conf"
	"github.com/justblue/luoye/services/goodbye/internal/domain"
	"github.com/justblue/luoye/services/goodbye/internal/event"
	"github.com/justblue/luoye/services/goodbye/internal/grpchandler"
	"github.com/justblue/luoye/services/goodbye/internal/server"
	"github.com/justblue/luoye/services/goodbye/internal/usecase"
	"github.com/nats-io/nats.go/jetstream"
)

func provideNATSClient(cfg *conf.Config) (*kitnats.Client, error) {
	client, err := kitnats.NewClient(cfg.NATS.URL)
	if err != nil {
		return nil, err
	}
	if err := client.EnsureStream(context.Background(), "goodbye", []string{"goodbye.said"}, jetstream.MemoryStorage); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

func provideAsynqClient(cfg *conf.Config) *asynq.Client {
	return kitasynq.NewClient(kitasynq.Config{Addr: cfg.Redis.Addr})
}

func initApp(cfg *conf.Config) (*kratos.App, func(), error) {
	wire.Build(
		provideNATSClient,
		provideAsynqClient,
		event.NewPublisher,
		event.NewAsynqPublisher,
		event.NewCompositePublisher,
		usecase.NewGoodbyeService,
		grpchandler.NewGoodbyeServer,
		server.NewGRPCServer,
		server.NewHTTPServer,
		server.NewApp,
		wire.Bind(new(domain.EventPublisher), new(*event.CompositePublisher)),
		wire.Bind(new(domain.Goodbye), new(*usecase.GoodbyeService)),
	)
	return nil, nil, nil
}
