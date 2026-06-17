package domain

import (
	"context"
	"time"
)

type GoodbyeSaidEvent struct {
	Name      string
	Message   string
	Timestamp time.Time
}

type GoodbyeTaskProcessor interface {
	ProcessGoodbyeSaid(ctx context.Context, evt GoodbyeSaidEvent) error
}
