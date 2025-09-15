package internal

import (
	"context"
	"errors"
	"testing"

	db "github.com/franzego/ecommerce-microservices/order-service/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrderWithListedItems_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderer)
	ctx := context.Background()
	userid := 1

	price := pgtype.Numeric{}
	price.Scan("99.99")
	price2 := pgtype.Numeric{}
	price2.Scan("209.99")
	// Mock data
	orderItems := []db.AddOrderItemParams{

		{
			ProductID: 1,
			Quantity:  2,
			Price:     price,
		},
		{
			ProductID: 2,
			Quantity:  1,
			Price:     price2,
		},
	}

	expectedOrder := db.Order{
		OrderID:     1,
		UserID:      1,
		TotalAmount: pgtype.Numeric{},
		StatusStaus: "pending",
		//CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
	}

	expectedOrderItems := []db.OrderItem{
		{
			OrderItemID: 1,
			OrderID:     1,
			ProductID:   1,
			Quantity:    2,
			Price:       price,
		},
		{
			OrderItemID: 2,
			OrderID:     1,
			ProductID:   2,
			Quantity:    1,
			Price:       price2,
		},
	}

	// Set expectations
	mockRepo.On("CreateOrderWithListedItems", ctx, userid, orderItems).Return(expectedOrder, expectedOrderItems, nil)

	// Act
	order, items, err := mockRepo.CreateOrderWithListedItems(ctx, userid, orderItems)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
	assert.Equal(t, expectedOrderItems, items)
	assert.Len(t, items, 2)
	mockRepo.AssertExpectations(t)
}

func TestCreateOrderWithListedItems_Error(t *testing.T) {

	mockRepo := new(MockOrderer)
	ctx := context.Background()
	userid := 1
	price := pgtype.Numeric{}
	price.Scan("99.99")
	orderItems := []db.AddOrderItemParams{
		{ProductID: 1, Quantity: 2, Price: price},
	}

	expectedError := errors.New("database connection failed")

	// Set expectations
	mockRepo.On("CreateOrderWithListedItems", ctx, userid, orderItems).Return(db.Order{}, []db.OrderItem(nil), expectedError)

	// Act
	order, items, err := mockRepo.CreateOrderWithListedItems(ctx, userid, orderItems)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, db.Order{}, order)
	assert.Nil(t, items)
	mockRepo.AssertExpectations(t)
}

func TestListOrders_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderer)
	ctx := context.Background()
	params := db.ListOrdersParams{
		//UserID: 1,
		Limit:  10,
		Offset: 0,
	}

	expectedOrders := []db.Order{
		{
			OrderID:     1,
			UserID:      1,
			StatusStaus: "pending",
			// CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			OrderID:     2,
			UserID:      1,
			StatusStaus: "completed",
			//CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	// Set expectations
	mockRepo.On("ListOrders", ctx, params).Return(expectedOrders, nil)

	// Act
	orders, err := mockRepo.ListOrders(ctx, params)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedOrders, orders)
	assert.Len(t, orders, 2)
	mockRepo.AssertExpectations(t)
}
