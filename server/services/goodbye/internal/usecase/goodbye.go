package usecase

import (
	"context"
	"fmt"

	"github.com/justblue/luoye/services/goodbye/internal/domain"
)

type GoodbyeService struct {
	publisher domain.EventPublisher
}

func NewGoodbyeService(publisher domain.EventPublisher) *GoodbyeService {
	return &GoodbyeService{publisher: publisher}
}

func (s *GoodbyeService) SayGoodbye(ctx context.Context, req *domain.GoodbyeRequest) (*domain.GoodbyeReply, error) {
	msg := fmt.Sprintf("Goodbye, %s!", req.Name)
	reply := &domain.GoodbyeReply{Message: msg}

	if err := s.publisher.PublishGoodbyeSaid(ctx, req.Name, msg); err != nil {
		return nil, err
	}

	return reply, nil
}
