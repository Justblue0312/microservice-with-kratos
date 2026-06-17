package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go/jetstream"

	kitnats "github.com/justblue/luoye/kit/messaging/nats"
)

func main() {
	url := flag.String("url", "nats://localhost:4222", "NATS server URL")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client, err := kitnats.NewClient(*url)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer client.Close()

	streams := map[string]struct {
		Subjects []string
		Storage  jetstream.StorageType
	}{
		"goodbye": {
			Subjects: []string{"goodbye.said"},
			Storage:  jetstream.MemoryStorage,
		},
	}

	for name, s := range streams {
		if err := client.EnsureStream(ctx, name, s.Subjects, s.Storage); err != nil {
			log.Fatalf("failed to ensure stream %q: %v", name, err)
		}
		log.Printf("stream %q ensured with subjects %v", name, s.Subjects)
	}

	log.Println("all streams provisioned successfully")
}
