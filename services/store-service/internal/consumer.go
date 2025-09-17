package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	srv    *Service
	reader *kafka.Reader
}

func NewConsumerService(brokers []string, srv *Service) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   "orders",
		GroupID: "store-service",
	})
	return &Consumer{
		srv:    srv,
		reader: reader,
	}
}

type Order struct {
	Productid int32 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

func (c *Consumer) HandleEvent(event []byte) {
	//var event []byte
	var parsed Order

	err := json.Unmarshal(event, &parsed)
	if err != nil {
		fmt.Printf("error in unmarshalling file: %v", err)
	}
	c.srv.StoreFunc(context.Background(), parsed.Productid, parsed.Quantity)

}

func (c *Consumer) ReadMessage(ctx context.Context) {
	defer c.reader.Close()
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error printing out message: %v", err)
			break
		}
		c.HandleEvent(msg.Value)
	}
}
