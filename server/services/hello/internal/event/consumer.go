package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	kitnats "github.com/justblue/luoye/kit/messaging/nats"
	"github.com/nats-io/nats.go/jetstream"
)

type Consumer struct {
	client *kitnats.Client
	mc     jetstream.MessagesContext
}

type goodbyeSaidEvent struct {
	Name      string `json:"name"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func NewConsumer(url string) (*Consumer, error) {
	client, err := kitnats.NewClient(url)
	if err != nil {
		return nil, err
	}
	if err := client.EnsureStream(context.Background(), "goodbye", []string{"goodbye.said"}, jetstream.MemoryStorage); err != nil {
		return nil, err
	}
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
	return &Consumer{client: client, mc: mc}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			msg, err := c.mc.Next()
			if err != nil {
				return
			}
			var evt goodbyeSaidEvent
			if err := json.Unmarshal(msg.Data(), &evt); err != nil {
				log.Printf("nats: bad message: %v", err)
				msg.Term()
				continue
			}
			log.Printf("nats: received goodbye.said: name=%s message=%q timestamp=%s",
				evt.Name, evt.Message, evt.Timestamp)
			fmt.Printf("Hello received via NATS: %s\n", evt.Message)
			msg.Ack()
		}
	}()
	return nil
}

func (c *Consumer) Stop(_ context.Context) error {
	if c.client != nil {
		c.client.Close()
	}
	return nil
}
