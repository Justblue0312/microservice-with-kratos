package server

import (
	"github.com/go-kratos/kratos/v2"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/justblue/luoye/services/hello/internal/event"
)

func NewApp(grpc *kratosgrpc.Server, consumer *event.Consumer) *kratos.App {
	return kratos.New(
		kratos.Name("hello"),
		kratos.Server(grpc, consumer),
	)
}
