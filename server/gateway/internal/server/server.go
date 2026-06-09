package server

import (
	"github.com/go-kratos/kratos/v2"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

func NewApp(http *kratoshttp.Server) *kratos.App {
	return kratos.New(
		kratos.Name("gateway"),
		kratos.Server(http),
	)
}
