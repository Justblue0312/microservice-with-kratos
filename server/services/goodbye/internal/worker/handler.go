package worker

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/hibiken/asynq"
	kitasynq "github.com/justblue/luoye/kit/messaging/asynq"
	"github.com/justblue/luoye/services/goodbye/internal/domain"
)

type GoodbyeHandler struct {
	processor domain.GoodbyeTaskProcessor
}

func NewGoodbyeHandler(processor domain.GoodbyeTaskProcessor) *GoodbyeHandler {
	return &GoodbyeHandler{processor: processor}
}

func (h *GoodbyeHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload kitasynq.GoodbyeSaidPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	evt := domain.GoodbyeSaidEvent{
		Name:    payload.Name,
		Message: payload.Message,
	}
	slog.InfoContext(ctx, "goodbye: received goodbye.said task",
		"name", payload.Name,
		"message", payload.Message,
	)
	return h.processor.ProcessGoodbyeSaid(ctx, evt)
}
