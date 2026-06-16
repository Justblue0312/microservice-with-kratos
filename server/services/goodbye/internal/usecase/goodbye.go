package usecase

import (
	"context"
	"fmt"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/justblue/luoye/services/goodbye/internal/domain"
)

type GoodbyeService struct {
	publisher domain.EventPublisher
}

func NewGoodbyeService(publisher domain.EventPublisher) *GoodbyeService {
	return &GoodbyeService{publisher: publisher}
}

func (s *GoodbyeService) SayGoodbye(ctx context.Context, req *domain.GoodbyeRequest) (*domain.GoodbyeReply, error) {
	if req.Name == "" {
		return nil, kerrors.BadRequest("GOODBYE_BAD_NAME", "name is required")
	}
	msg := fmt.Sprintf("Goodbye, %s!", req.Name)
	reply := &domain.GoodbyeReply{Message: msg}

	if err := s.publisher.PublishGoodbyeSaid(ctx, req.Name, msg); err != nil {
		return nil, err
	}

	return reply, nil
}
