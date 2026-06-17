package server

import (
	"context"
	"sync/atomic"

	"github.com/hibiken/asynq"
	kitasynqtasks "github.com/justblue/luoye/kit/messaging/asynq"
	"github.com/justblue/luoye/services/goodbye/internal/worker"
)

type AsynqService struct {
	server  *asynq.Server
	mux     *asynq.ServeMux
	started atomic.Bool
}

func NewAsynqService(server *asynq.Server, h *worker.GoodbyeHandler) *AsynqService {
	mux := asynq.NewServeMux()
	mux.Handle(kitasynqtasks.TypeGoodbyeSaid, h)
	return &AsynqService{server: server, mux: mux}
}

func (s *AsynqService) Start(ctx context.Context) error {
	if err := s.server.Start(s.mux); err != nil {
		return err
	}
	s.started.Store(true)
	<-ctx.Done()
	return nil
}

func (s *AsynqService) Stop(_ context.Context) error {
	if s.started.Load() {
		s.server.Shutdown()
	}
	return nil
}
