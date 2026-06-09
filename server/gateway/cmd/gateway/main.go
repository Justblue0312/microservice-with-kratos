package main

import (
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/justblue/luoye/gateway/internal/conf"
	"github.com/justblue/luoye/gateway/internal/proxy"
	"github.com/justblue/luoye/gateway/internal/server"
)

func initApp(cfg *conf.Config) (*kratos.App, error) {
	helloProxy, err := proxy.NewHelloProxy(cfg.Upstreams.Hello)
	if err != nil {
		return nil, err
	}
	httpServer := server.NewHTTPServer(cfg, helloProxy)
	return server.NewApp(httpServer), nil
}

func main() {
	cfg := &conf.Config{
		HTTP:      conf.HTTPConfig{Addr: ":8080"},
		Upstreams: conf.UpstreamsConfig{Hello: "localhost:9081"},
	}
	app, err := initApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
