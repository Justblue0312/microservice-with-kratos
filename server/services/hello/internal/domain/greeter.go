package domain

import "context"

type GreetRequest struct {
    Name string
}

type GreetReply struct {
    Message string
}

// GreeterService is the interface the app layer implements.
// Repository interfaces go here too when you have DB access.
type Greeter interface {
    Greet(ctx context.Context, req *GreetRequest) (*GreetReply, error)
}
