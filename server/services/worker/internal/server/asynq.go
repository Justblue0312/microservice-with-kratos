package server

import (
	"context"

	"github.com/hibiken/asynq"
	kitasynqtasks "github.com/justblue/luoye/kit/messaging/asynq"
	"github.com/justblue/luoye/services/worker/internal/handler"
)

type AsynqService struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewAsynqService(server *asynq.Server, h *handler.GoodbyeHandler) *AsynqService {
	mux := asynq.NewServeMux()
	mux.Handle(kitasynqtasks.TypeGoodbyeSaid, h)
	return &AsynqService{server: server, mux: mux}
}

func (s *AsynqService) Start(ctx context.Context) error {
	if err := s.server.Start(s.mux); err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *AsynqService) Stop(ctx context.Context) error {
	s.server.Shutdown()
	return nil
}
