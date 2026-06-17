package nats

import (
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Client struct {
	conn *nats.Conn
	js   jetstream.JetStream
}

func NewClient(url string) (*Client, error) {
	nc, err := nats.Connect(url,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			slog.Error("nats: disconnected", "err", err)
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			slog.Info("nats: reconnected")
		}),
	)
	if err != nil {
		return nil, err
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	return &Client{conn: nc, js: js}, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Drain()
	}
}

func (c *Client) Conn() *nats.Conn        { return c.conn }
func (c *Client) JS() jetstream.JetStream { return c.js }
