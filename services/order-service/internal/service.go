package internal

import (
	"context"
	"log"

	db "github.com/franzego/ecommerce-microservices/order-service/db/sqlc"
)

type Servicer interface {
	CreateOrderWithItems(ctx context.Context, userid int, args []db.AddOrderItemParams) (db.Order, []db.OrderItem, error)
	ListOrders(ctx context.Context, args db.ListOrdersParams) ([]db.Order, error)
	UpdateOrderStatus(ctx context.Context, orderid int32, statusstatus string) (*db.Order, error)
	GetOrder(ctx context.Context, orderid int32) (*db.Order, error)
}

type Service struct {
	repo Orderer
}

func NewService(ord Orderer) *Service {
	if ord == nil {
		return nil
	}
	return &Service{repo: ord}
}

// func to actually create the order
func (s *Service) CreateOrderWithItems(ctx context.Context, userid int, args []db.AddOrderItemParams) (db.Order, []db.OrderItem, error) {
	ord, orditems, err := s.repo.CreateOrderWithListedItems(ctx, userid, args)
	if err != nil {
		log.Printf("error in CreatingOrder service: %v", err)
		return db.Order{}, nil, err
	}
	return ord, orditems, nil
}

// func to actually list order
func (s *Service) ListOrders(ctx context.Context, args db.ListOrdersParams) ([]db.Order, error) {
	ordlist, err := s.repo.ListOrders(ctx, args)
	if err != nil {
		log.Printf("error in ListOrder service: %v", err)
		return nil, err
	}
	return ordlist, nil
}

// func to actually update the order status
func (s *Service) UpdateOrderStatus(ctx context.Context, orderid int32, statusstatus string) (*db.Order, error) {
	updatedorder, err := s.repo.UpdateOrderStatus(ctx, orderid, statusstatus)
	if err != nil {
		log.Printf("error in UpdateOrder service: %v", err)
		return nil, err
	}
	return &updatedorder, nil
}

// func to get an order
func (s *Service) GetOrder(ctx context.Context, orderid int32) (*db.Order, error) {
	gottenorder, err := s.repo.GetOrder(ctx, orderid)
	if err != nil {
		log.Printf("error in GetOrder service: %v", err)
		return nil, err
	}
	return &gottenorder, nil
}
