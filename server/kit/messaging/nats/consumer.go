package nats

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type Subscription struct {
	consumer jetstream.Consumer
}

func (c *Client) Subscribe(ctx context.Context, stream string, cfg jetstream.ConsumerConfig) (*Subscription, error) {
	consumer, err := c.js.CreateOrUpdateConsumer(ctx, stream, cfg)
	if err != nil {
		return nil, err
	}
	return &Subscription{consumer: consumer}, nil
}

func (s *Subscription) Messages() (jetstream.MessagesContext, error) {
	return s.consumer.Messages()
}
