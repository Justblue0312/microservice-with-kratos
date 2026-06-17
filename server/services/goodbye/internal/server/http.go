package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	goodbyev1 "github.com/justblue/luoye/gen/go/goodbye"
	"github.com/justblue/luoye/services/goodbye/internal/conf"
	grpchandler "github.com/justblue/luoye/services/goodbye/internal/grpchandler"
)

func NewHTTPServer(cfg *conf.Config, goodbye *grpchandler.GoodbyeServer) *kratoshttp.Server {
	srv := kratoshttp.NewServer(
		kratoshttp.Address(cfg.HTTP.Addr),
		kratoshttp.Middleware(
			recovery.Recovery(),
		),
	)

	redisOpt := asynq.RedisClientOpt{Addr: cfg.Redis.Addr}
	inspector := asynq.NewInspector(redisOpt)
	mon := asynqmon.New(asynqmon.Options{
		RootPath:     "/monitor",
		RedisConnOpt: redisOpt,
	})
	srv.HandlePrefix("/monitor", mon)

	r := srv.Route("/")
	r.GET("/health", func(ctx kratoshttp.Context) error {
		return ctx.Result(http.StatusOK, map[string]string{"status": "ok"})
	})
	r.GET("/metrics", func(ctx kratoshttp.Context) error {
		queues, err := inspector.Queues()
		if err != nil {
			return ctx.Result(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		type queueStat struct {
			Queue   string `json:"queue"`
			Active  int    `json:"active"`
			Pending int    `json:"pending"`
			Failed  int    `json:"failed"`
			Size    int    `json:"size"`
		}
		var stats []queueStat
		for _, q := range queues {
			info, err := inspector.GetQueueInfo(q)
			if err != nil {
				continue
			}
			stats = append(stats, queueStat{
				Queue:   q,
				Active:  info.Active,
				Pending: info.Pending,
				Failed:  info.Failed,
				Size:    info.Size,
			})
		}
		w := ctx.Response()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(stats)
		return nil
	})
	goodbyev1.RegisterGoodbyeHTTPServer(srv, goodbye)
	return srv
}
