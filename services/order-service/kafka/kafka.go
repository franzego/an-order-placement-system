package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/franzego/ecommerce-microservices/order-service/events"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishEvents(ctx context.Context, topic string, event events.Event) error
	Close() error
}

type KafkaProducer struct {
	writers map[string]*kafka.Writer
	enabled bool
	brokers []string
}

func NewKafkaProducer(brokers []string, enabled bool) *KafkaProducer {
	if !enabled {
		log.Println("Kafka is disabled")
		return &KafkaProducer{
			writers: make(map[string]*kafka.Writer),
			enabled: false,
			brokers: brokers,
		}
	}

	// Check if Kafka is available
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		log.Printf("Warning: Could not connect to Kafka: %v", err)
		return &KafkaProducer{
			writers: make(map[string]*kafka.Writer),
			enabled: false,
			brokers: brokers,
		}
	}
	defer conn.Close()

	return &KafkaProducer{
		writers: make(map[string]*kafka.Writer),
		enabled: enabled,
		brokers: brokers,
	}
}
func (kp *KafkaProducer) getorCreateWriter(topic string, brokers []string) *kafka.Writer {
	if writer, exists := kp.writers[topic]; exists {
		return writer
	}
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireOne,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
	}

	kp.writers[topic] = writer
	return writer
}
func (kp *KafkaProducer) PublishEvents(ctx context.Context, topic string, event events.Event) error {
	if !kp.enabled {
		log.Printf("Kafka is disabled, skipping event publishing for topic %s", topic)
		return nil
	}

	// Check topic exists
	conn, err := kafka.Dial("tcp", kp.brokers[0])
	if err != nil {
		log.Printf("Warning: Could not connect to Kafka: %v", err)
		return nil
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		log.Printf("Warning: Topic %s does not exist: %v", topic, err)
		return fmt.Errorf("topic %s does not exist: %v", topic, err)
	}
	if len(partitions) == 0 {
		log.Printf("Warning: Topic %s has no partitions", topic)
		return fmt.Errorf("topic %s has no partitions", topic)
	}

	evt, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error while marshalling %v", err)
	}
	writer := kp.getorCreateWriter(topic, kp.brokers)

	message := kafka.Message{
		Key:   []byte(event.ID), // Use event ID as message key for ordering
		Value: evt,
		Headers: []kafka.Header{
			{Key: "event-type", Value: []byte(event.Type)},
			{Key: "event-version", Value: []byte(event.Version)},
		},
	}
	if err := writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to write message to topic %s: %w", topic, err)
	}
	log.Printf("Published event %s to topic %s", event.Type, topic)
	return nil
}
func (kp *KafkaProducer) Close() error {
	for topic, writer := range kp.writers {
		if err := writer.Close(); err != nil {
			log.Printf("Error closing writer for topic %s: %v", topic, err)
		}
	}
	return nil
}
