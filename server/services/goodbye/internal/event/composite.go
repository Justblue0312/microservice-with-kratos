package event

import "context"

type CompositePublisher struct {
	nats  *Publisher
	asynq *AsynqPublisher
}

func NewCompositePublisher(nats *Publisher, asynq *AsynqPublisher) *CompositePublisher {
	return &CompositePublisher{nats: nats, asynq: asynq}
}

func (p *CompositePublisher) PublishGoodbyeSaid(ctx context.Context, name, message string) error {
	if err := p.nats.PublishGoodbyeSaid(ctx, name, message); err != nil {
		return err
	}
	return p.asynq.PublishGoodbyeSaid(ctx, name, message)
}
