//go:build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	"github.com/hibiken/asynq"
	kitasynq "github.com/justblue/luoye/kit/messaging/asynq"
	"github.com/justblue/luoye/services/worker/internal/conf"
	"github.com/justblue/luoye/services/worker/internal/handler"
	"github.com/justblue/luoye/services/worker/internal/server"
)

func provideAsynqServer(cfg *conf.Config) *asynq.Server {
	return kitasynq.NewServer(kitasynq.Config{Addr: cfg.Redis.Addr})
}

func initApp(cfg *conf.Config) (*kratos.App, func(), error) {
	wire.Build(
		provideAsynqServer,
		handler.NewGoodbyeHandler,
		server.NewHTTPServer,
		server.NewAsynqService,
		server.NewApp,
	)
	return nil, nil, nil
}
