package usecase

import (
	"context"
	"fmt"

	"github.com/justblue/luoye/services/hello/internal/domain"
)

type GreeterService struct {
	goodbye domain.GoodbyeClient
}

func NewGreeterService(goodbye domain.GoodbyeClient) *GreeterService {
	return &GreeterService{goodbye: goodbye}
}

func (s *GreeterService) Greet(ctx context.Context, req *domain.GreetRequest) (*domain.GreetReply, error) {
	greeting := fmt.Sprintf("Hello, %s!", req.Name)

	goodbyeReply, err := s.goodbye.SayGoodbye(ctx, &domain.GoodbyeRequest{Name: req.Name})
	if err != nil {
		return nil, err
	}

	return &domain.GreetReply{
		Message: fmt.Sprintf("%s %s", greeting, goodbyeReply.Message),
	}, nil
}
