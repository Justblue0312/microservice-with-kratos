package domain

import "context"

type GreetRequest struct {
	Name string
}

type GreetReply struct {
	Message string
}

type GoodbyeRequest struct {
	Name string
}

type GoodbyeReply struct {
	Message string
}

type Greeter interface {
	Greet(ctx context.Context, req *GreetRequest) (*GreetReply, error)
}

type GoodbyeClient interface {
	SayGoodbye(ctx context.Context, req *GoodbyeRequest) (*GoodbyeReply, error)
}
