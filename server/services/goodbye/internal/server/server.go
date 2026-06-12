package server

import (
	"github.com/go-kratos/kratos/v2"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
)

func NewApp(grpc *kratosgrpc.Server) *kratos.App {
	return kratos.New(
		kratos.Name("goodbye"),
		kratos.Server(grpc),
	)
}
