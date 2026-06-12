package event

import (
	"context"
	"encoding/json"
	"time"

	kitnats "github.com/justblue/luoye/kit/messaging/nats"
	"github.com/nats-io/nats.go/jetstream"
)

type Publisher struct {
	client *kitnats.Client
}

type goodbyeSaidEvent struct {
	Name      string `json:"name"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func NewPublisher(url string) (*Publisher, error) {
	client, err := kitnats.NewClient(url)
	if err != nil {
		return nil, err
	}
	if err := client.EnsureStream(context.Background(), "goodbye", []string{"goodbye.said"}, jetstream.MemoryStorage); err != nil {
		return nil, err
	}
	return &Publisher{client: client}, nil
}

func (p *Publisher) PublishGoodbyeSaid(_ context.Context, name, message string) error {
	data, err := json.Marshal(goodbyeSaidEvent{
		Name:      name,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	return p.client.Publish(context.Background(), "goodbye.said", data)
}
