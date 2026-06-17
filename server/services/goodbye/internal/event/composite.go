package event

import (
	"context"
	"log/slog"
)

type CompositePublisher struct {
	nats  *Publisher
	asynq *AsynqPublisher
}

func NewCompositePublisher(nats *Publisher, asynq *AsynqPublisher) *CompositePublisher {
	return &CompositePublisher{nats: nats, asynq: asynq}
}

func (p *CompositePublisher) PublishGoodbyeSaid(ctx context.Context, name, message string) error {
	if err := p.nats.PublishGoodbyeSaid(ctx, name, message); err != nil {
		slog.ErrorContext(ctx, "goodbye: failed to publish to NATS",
			"name", name, "err", err,
		)
	}
	if err := p.asynq.PublishGoodbyeSaid(ctx, name, message); err != nil {
		slog.ErrorContext(ctx, "goodbye: failed to enqueue to Asynq",
			"name", name, "err", err,
		)
		return err
	}
	return nil
}
