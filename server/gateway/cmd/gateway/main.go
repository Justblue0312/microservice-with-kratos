package main

import (
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/justblue/luoye/gateway/internal/conf"
	"github.com/justblue/luoye/gateway/internal/server"
)

func initApp(cfg *conf.Config) (*kratos.App, error) {
	httpServer, err := server.NewHTTPServer(cfg)
	if err != nil {
		return nil, err
	}
	return server.NewApp(httpServer), nil
}

func main() {
	c := config.New(
		config.WithSource(file.NewSource("gateway/configs/config.yaml")),
	)
	if err := c.Load(); err != nil {
		log.Fatal(err)
	}
	var cfg conf.Config
	if err := c.Scan(&cfg); err != nil {
		log.Fatal(err)
	}
	app, err := initApp(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
