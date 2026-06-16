package main

import (
	"log"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/justblue/luoye/services/hello/internal/conf"
)

func main() {
	c := config.New(
		config.WithSource(file.NewSource("services/hello/config/config.yaml")),
	)
	if err := c.Load(); err != nil {
		log.Fatal(err)
	}
	var cfg conf.Config
	if err := c.Scan(&cfg); err != nil {
		log.Fatal(err)
	}
	app, cleanup, err := initApp(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
