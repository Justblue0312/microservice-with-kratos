package server

import (
	"github.com/go-kratos/kratos/v2"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/justblue/luoye/services/hello/internal/event"
)

func NewApp(http *kratoshttp.Server, grpc *kratosgrpc.Server, consumer *event.Consumer) *kratos.App {
	return kratos.New(
		kratos.Name("hello"),
		kratos.Server(http, grpc, consumer),
	)
}
