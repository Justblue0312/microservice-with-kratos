package nats

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

func (c *Client) EnsureStream(ctx context.Context, name string, subjects []string, storage jetstream.StorageType) error {
	_, err := c.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     name,
		Subjects: subjects,
		Storage:  storage,
	})
	return err
}

func (c *Client) Publish(ctx context.Context, subject string, data []byte) error {
	_, err := c.js.Publish(ctx, subject, data)
	return err
}
