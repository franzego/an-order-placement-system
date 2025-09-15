package internal

import (
	"context"

	db "github.com/franzego/ecommerce-microservices/order-service/db/sqlc"
	"github.com/stretchr/testify/mock"
)

// MockOrderer is a mock implementation of the Orderer interface
type MockOrderer struct {
	mock.Mock
}

func (m *MockOrderer) CreateOrderWithListedItems(ctx context.Context, userid int, args []db.AddOrderItemParams) (db.Order, []db.OrderItem, error) {
	ret := m.Called(ctx, userid, args)
	return ret.Get(0).(db.Order), ret.Get(1).([]db.OrderItem), ret.Error(2)
}

func (m *MockOrderer) ListOrders(ctx context.Context, args db.ListOrdersParams) ([]db.Order, error) {
	ret := m.Called(ctx, args)
	return ret.Get(0).([]db.Order), ret.Error(1)
}

func (m *MockOrderer) UpdateOrderStatus(ctx context.Context, orderid int32, statusstatus string) (db.Order, error) {
	ret := m.Called(ctx, orderid, statusstatus)
	return ret.Get(0).(db.Order), ret.Error(1)
}

func (m *MockOrderer) GetOrder(ctx context.Context, orderid int32) (db.Order, error) {
	ret := m.Called(ctx, orderid)
	return ret.Get(0).(db.Order), ret.Error(1)
}
