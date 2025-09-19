package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/franzego/store-service/config"
	"github.com/franzego/store-service/internal"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.LoadConfig()
	dbconn2, err := pgxpool.New(context.Background(), cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer dbconn2.Close()

	repo := internal.NewRepoService(dbconn2)
	if repo == nil {
		log.Fatal("Failed to create repository service")
	}

	svc := internal.NewService(repo)
	if svc == nil {
		log.Fatal("Failed to create service")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handle shutdown signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ch
		log.Println("shutting down store service...")
		cancel()
	}()

	if cfg.KafkaEnabled {
		cons := internal.NewConsumerService(cfg.KafkaBrokers, svc)
		if cons == nil {
			log.Fatal("Failed to initiate consumer service")
		}
		cons.ReadMessage(ctx)
	} else {
		log.Println("Kafka consumer disabled by config; waiting for shutdown signal")
		<-ctx.Done()
	}
}
