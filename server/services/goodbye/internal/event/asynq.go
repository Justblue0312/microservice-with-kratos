package event

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	kitasynq "github.com/justblue/luoye/kit/messaging/asynq"
)

type AsynqPublisher struct {
	client *asynq.Client
}

func NewAsynqPublisher(client *asynq.Client) *AsynqPublisher {
	return &AsynqPublisher{client: client}
}

func (p *AsynqPublisher) PublishGoodbyeSaid(ctx context.Context, name, message string) error {
	task, err := kitasynq.NewGoodbyeSaidTask(name, message, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return err
	}
	_, err = p.client.EnqueueContext(ctx, task)
	return err
}
