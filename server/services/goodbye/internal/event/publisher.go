package event

import (
	"context"
	"encoding/json"
	"time"

	kitnats "github.com/justblue/luoye/kit/messaging/nats"
)

type Publisher struct {
	client *kitnats.Client
}

type goodbyeSaidEvent struct {
	Name      string `json:"name"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func NewPublisher(client *kitnats.Client) (*Publisher, error) {
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
