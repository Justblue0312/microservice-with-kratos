package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
	kitasynq "github.com/justblue/luoye/kit/messaging/asynq"
)

var _ asynq.Handler = (*GoodbyeHandler)(nil)

type GoodbyeHandler struct{}

func NewGoodbyeHandler() *GoodbyeHandler {
	return &GoodbyeHandler{}
}

func (h *GoodbyeHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload kitasynq.GoodbyeSaidPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	log.Printf("worker: processed goodbye.said event — name=%s message=%q timestamp=%s",
		payload.Name, payload.Message, payload.Timestamp)
	return nil
}
