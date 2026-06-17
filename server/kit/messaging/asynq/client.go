package asynq

import (
	"context"
	"log/slog"

	"github.com/hibiken/asynq"
)

type Config struct {
	Addr        string
	Concurrency int
	Queues      map[string]int
}

func NewClient(cfg Config) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.Addr})
}

func NewServer(cfg Config) *asynq.Server {
	concurrency := cfg.Concurrency
	if concurrency == 0 {
		concurrency = 10
	}
	queues := cfg.Queues
	if len(queues) == 0 {
		queues = map[string]int{"default": 1}
	}
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Addr},
		asynq.Config{
			Concurrency: concurrency,
			Queues:      queues,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				retried, _ := asynq.GetRetryCount(ctx)
				maxRetry, _ := asynq.GetMaxRetry(ctx)
				slog.ErrorContext(ctx, "asynq: task failed",
					"type", task.Type(),
					"retry", retried,
					"max_retry", maxRetry,
					"err", err,
				)
			}),
		},
	)
}
