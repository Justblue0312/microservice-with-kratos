package domain

import "context"

type GoodbyeRequest struct {
	Name string
}

type GoodbyeReply struct {
	Message string
}

type Goodbye interface {
	SayGoodbye(ctx context.Context, req *GoodbyeRequest) (*GoodbyeReply, error)
}

type EventPublisher interface {
	PublishGoodbyeSaid(ctx context.Context, name, message string) error
}
