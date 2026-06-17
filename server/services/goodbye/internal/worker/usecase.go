package worker

import (
	"context"
	"log/slog"

	"github.com/justblue/luoye/services/goodbye/internal/domain"
)

type GoodbyeProcessor struct {
}

func NewGoodbyeProcessor() *GoodbyeProcessor {
	return &GoodbyeProcessor{}
}

func (p *GoodbyeProcessor) ProcessGoodbyeSaid(ctx context.Context, evt domain.GoodbyeSaidEvent) error {
	slog.InfoContext(ctx, "goodbye: processing goodbye.said event",
		"name", evt.Name,
		"message", evt.Message,
		"timestamp", evt.Timestamp,
	)
	return nil
}
