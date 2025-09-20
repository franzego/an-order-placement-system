package main

import (
	"context"
	"log"

	"github.com/franzego/notification-service/internal"
)

func main() {
	// Redis configuration
	redisConfig := internal.RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		//PoolSize: 10,
	}

	consumer, err := internal.NewConsumerService(
		[]string{"localhost:9092"}, // Kafka brokers
		redisConfig,
	)
	if err != nil {
		log.Fatalf("Failed to initialize consumer service: %v", err)
	}

	// Start consuming messages
	go consumer.ReadMessage(context.Background())

	// Keep the application running
	<-context.Background().Done()
}
