package internal

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

/*
	func StartConsumer(brokers []string, topic, groupID string, handler func(ctx context.Context, msg []byte) error) {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		})
		defer r.Close()

		ctx := context.Background()
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				log.Printf(" error reading message: %v", err)
				continue
			}

			if err := handler(ctx, m.Value); err != nil {
				log.Printf("error handling message: %v", err)
			}
		}
	}
*/
type Consumer struct {
	not    *Notification
	reader *kafka.Reader
}

func NewConsumerService(brokers []string, redisConfig RedisConfig) (*Consumer, error) {
	// Initialize notification service with Redis
	notificationService, err := NewNotificationService(redisConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize notification service: %w", err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          "orders",
		GroupID:        "notification-service",
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		ReadBackoffMin: time.Second,
		ReadBackoffMax: 5 * time.Second,
	})

	return &Consumer{
		not:    notificationService,
		reader: reader,
	}, nil
}

func (c *Consumer) ReadMessage(ctx context.Context) {
	defer c.reader.Close()
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf(" error reading message: %v", err)
			continue
		}
		log.Printf("consumed message from topic %s, partition %d, offset %d", msg.Topic, msg.Partition, msg.Offset)
		c.not.HandleEvent(ctx, msg.Value)
	}

}
