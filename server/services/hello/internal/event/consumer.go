package event

import (
	"context"
	"encoding/json"
	"log/slog"

	kitnats "github.com/justblue/luoye/kit/messaging/nats"
	"github.com/nats-io/nats.go/jetstream"
)

type Consumer struct {
	mc jetstream.MessagesContext
}

type goodbyeSaidEvent struct {
	Name      string `json:"name"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func NewConsumer(client *kitnats.Client) (*Consumer, error) {
	sub, err := client.Subscribe(context.Background(), "goodbye", jetstream.ConsumerConfig{
		Name:          "hello-goodbye-consumer",
		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return nil, err
	}
	mc, err := sub.Messages()
	if err != nil {
		return nil, err
	}
	return &Consumer{mc: mc}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	go func() {
		for {
			msg, err := c.mc.Next()
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				slog.ErrorContext(ctx, "nats: consumer error, stopping",
					"err", err,
				)
				return
			}
			var evt goodbyeSaidEvent
			if err := json.Unmarshal(msg.Data(), &evt); err != nil {
				slog.ErrorContext(ctx, "nats: bad message",
					"err", err,
				)
				msg.Term()
				continue
			}
			slog.InfoContext(ctx, "nats: received goodbye.said",
				"name", evt.Name,
				"message", evt.Message,
				"timestamp", evt.Timestamp,
			)
			msg.Ack()
		}
	}()
	return nil
}

func (c *Consumer) Stop(_ context.Context) error {
	if c.mc != nil {
		c.mc.Stop()
	}
	return nil
}
