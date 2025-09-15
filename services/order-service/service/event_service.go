package service

import (
	"context"
	"fmt"
	"strconv"

	db "github.com/franzego/ecommerce-microservices/order-service/db/sqlc"
	"github.com/franzego/ecommerce-microservices/order-service/events"
	"github.com/franzego/ecommerce-microservices/order-service/kafka"
)

type EventService struct {
	Producer kafka.Producer
}

func NewEventService(producer kafka.Producer) *EventService {
	return &EventService{Producer: producer}
}

const (
	OrdersTopic = "orders"
)

func (es *EventService) PublishOrderCreated(ctx context.Context, order db.Order, items []db.OrderItem) error {
	// Convert order items to event format
	eventItems := make([]events.OrderItemData, len(items))
	for i, item := range items {
		priceVal, err := item.Price.Value()
		if err != nil {
			return fmt.Errorf("failed to get item price value: %w", err)
		}

		priceStr, ok := priceVal.(string)
		if !ok {
			return fmt.Errorf("expected string from numeric value, got %T", priceVal)
		}

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return fmt.Errorf("failed to parse price %s: %w", priceStr, err)
		}

		eventItems[i] = events.OrderItemData{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
		}
	}

	// Use Value() method for total amount
	totalVal, err := order.TotalAmount.Value()
	if err != nil {
		return fmt.Errorf("failed to get total amount value: %w", err)
	}

	totalAmount, ok := totalVal.(string)
	if !ok {
		return fmt.Errorf("expected string from numeric value, got %T", totalVal)
	}

	event := events.NewOrderCreatedService(
		int(order.OrderID),
		int(order.UserID),
		totalAmount,
		eventItems,
	)

	return es.Producer.PublishEvents(ctx, OrdersTopic, *event)
}

func (es *EventService) PublishOrderStatusUpdated(ctx context.Context, orderID, userID int32, oldStatus, newStatus string) error {
	event := events.NewUpdatedOrderService(int(orderID), oldStatus, newStatus, int(userID))
	return es.Producer.PublishEvents(ctx, OrdersTopic, *event)
}
