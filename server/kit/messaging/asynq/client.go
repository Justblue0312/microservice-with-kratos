package asynq

import "github.com/hibiken/asynq"

type Config struct {
	Addr string
}

func NewClient(cfg Config) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.Addr})
}

func NewServer(cfg Config) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Addr},
		asynq.Config{
			Concurrency: 10,
			Queues:      map[string]int{"default": 1},
		},
	)
}
