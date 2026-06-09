package usecase

import (
	"context"
	"fmt"

	"github.com/justblue/luoye/services/hello/internal/domain"
)

type GreeterService struct{}

func NewGreeterService() *GreeterService {
    return &GreeterService{}
}

func (s *GreeterService) Greet(_ context.Context, req *domain.GreetRequest) (*domain.GreetReply, error) {
    return &domain.GreetReply{Message: fmt.Sprintf("Hello, %s!", req.Name)}, nil
}
