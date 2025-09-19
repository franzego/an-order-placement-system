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

type OrderEvent struct {
	Data OrderData `json:"data"`
}

type OrderData struct {
	Items []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductID int32 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

/*type Order struct {
	Productid int32 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}*/

func (c *Consumer) HandleEvent(event []byte) {
	//var event []byte
	//var parsed Order
	var parsed OrderEvent

	err := json.Unmarshal(event, &parsed)
	if err != nil {
		fmt.Printf("error in unmarshalling file: %v", err)
	}
	//c.srv.StoreFunc(context.Background(), parsed.Productid, parsed.Quantity)
	for _, item := range parsed.Data.Items {
		c.srv.StoreFunc(context.Background(), item.ProductID, item.Quantity)
	}

}

func (c *Consumer) ReadMessage(ctx context.Context) {
	defer c.reader.Close()
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error printing out message: %v", err)
			break
		}
		log.Printf("consumed message from topic %s, partition %d, offset %d", msg.Topic, msg.Partition, msg.Offset)
		c.HandleEvent(msg.Value)

	}

}
