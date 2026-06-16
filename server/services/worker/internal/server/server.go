package server

import (
	"github.com/go-kratos/kratos/v2"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

func NewApp(hs *kratoshttp.Server, as *AsynqService) *kratos.App {
	return kratos.New(
		kratos.Name("worker"),
		kratos.Server(hs, as),
	)
}
