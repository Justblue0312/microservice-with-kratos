//go:build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	kitnats "github.com/justblue/luoye/kit/messaging/nats"
	"github.com/justblue/luoye/services/hello/internal/conf"
	"github.com/justblue/luoye/services/hello/internal/domain"
	"github.com/justblue/luoye/services/hello/internal/event"
	"github.com/justblue/luoye/services/hello/internal/grpcclient"
	"github.com/justblue/luoye/services/hello/internal/grpchandler"
	"github.com/justblue/luoye/services/hello/internal/server"
	"github.com/justblue/luoye/services/hello/internal/usecase"
)

func provideNATSClient(cfg *conf.Config) (*kitnats.Client, func(), error) {
	client, err := kitnats.NewClient(cfg.NATS.URL)
	if err != nil {
		return nil, nil, err
	}
	return client, func() { client.Close() }, nil
}

func provideUpstream(cfg *conf.Config) conf.UpstreamConfig { return cfg.Upstream }

func initApp(cfg *conf.Config) (*kratos.App, func(), error) {
	wire.Build(
		provideNATSClient,
		provideUpstream,
		grpcclient.NewGoodbyeClient,
		usecase.NewGreeterService,
		grpchandler.NewGreeterServer,
		server.NewGRPCServer,
		server.NewHTTPServer,
		event.NewConsumer,
		server.NewApp,
		wire.Bind(new(domain.Greeter), new(*usecase.GreeterService)),
		wire.Bind(new(domain.GoodbyeClient), new(*grpcclient.GoodbyeClient)),
	)
	return nil, nil, nil
}
