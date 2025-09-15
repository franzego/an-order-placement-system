package events

import (
	"fmt"
	"time"
)

type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Version   string      `json:"version"`
	Data      interface{} `json:"data"`
}

//ordercreated event

type OrderCreatedEvent struct {
	OrderID     int              `json:"orderid"`
	UserID      int              `json:"userid"`
	TotalAmount string           `json:"totalamount"`
	Status      string           `json:"status"`
	Items       []OrderItemData  `json:"items"`
}

type UpdatedOrderEvent struct {
	OrderID   int    `json:"orderid"`
	OldStatus string `json:"oldstatus"`
	NewStatus string `json:"newstatus"`
	UserID    int    `json:"userid"`
}

type OrderItemData struct {
	ProductID int32   `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Price     float64 `json:"price"`
}

const (
	OrderCreatedEventType       = "order.created"
	OrderStatusUpdatedEventType = "order.status_updated"
)

func generateEventID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// we need methods to initiate this structs
func NewOrderCreatedService(orderid int, userid int, totalamount string, items []OrderItemData) *Event {
	return &Event{
		ID:        generateEventID(),
		Type:      OrderCreatedEventType,
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
		Data: OrderCreatedEvent{
			OrderID:     orderid,
			UserID:      userid,
			TotalAmount: totalamount,
			Items:       items,
			Status:      "pending",
		},
	}
}
func NewUpdatedOrderService(orderid int, oldstatus string, newstatus string, userid int) *Event {
	return &Event{
		ID:        generateEventID(),
		Type:      OrderStatusUpdatedEventType,
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
		Data: UpdatedOrderEvent{
			OrderID:   orderid,
			OldStatus: oldstatus,
			NewStatus: newstatus,
			UserID:    userid,
		},
	}
}
