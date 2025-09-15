package internal

import (
	"context"
	"fmt"
	"log"
	"strconv"

	db "github.com/franzego/ecommerce-microservices/order-service/db/sqlc"
	"github.com/franzego/ecommerce-microservices/order-service/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Orderer interface {
	CreateOrderWithListedItems(ctx context.Context, userid int, args []db.AddOrderItemParams) (db.Order, []db.OrderItem, error)
	ListOrders(ctx context.Context, args db.ListOrdersParams) ([]db.Order, error)
	UpdateOrderStatus(ctx context.Context, orderid int32, statusstatus string) (db.Order, error)
	GetOrder(ctx context.Context, orderid int32) (db.Order, error)
}

type Repo struct {
	q            *db.Queries
	db           *pgxpool.Pool
	eventservice *service.EventService
}

func NewRepoService(dbConn *pgxpool.Pool, eventservice *service.EventService) *Repo {
	if dbConn == nil {
		return nil
	}

	repo := &Repo{
		q:  db.New(dbConn),
		db: dbConn,
	}

	if eventservice != nil {
		repo.eventservice = eventservice
	}

	return repo
}

//func NewRepoService(store OrderMock) *Repo {
//	return &Repo{q: store}
//}

// function to initiate the transaction
func (r *Repo) execTx(ctx context.Context, f func(*db.Queries) error) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error in initiating transactions %v", err)
	}
	qtx := db.New(tx)
	err = f(qtx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, tb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

// ///
func FloattoNumeric(f float64) pgtype.Numeric {
	var num pgtype.Numeric
	num.Valid = true // Ensure the numeric is valid
	// Format with exactly 2 decimal places and scan
	str := fmt.Sprintf("%.2f", f)
	log.Printf("Converting float %v to numeric string: %s", f, str)
	err := num.Scan(str)
	if err != nil {
		log.Printf("failed to convert %f to numeric: %v", f, err)
		num.Valid = true
		num.Scan("0.00")
	}
	return num
}

// ///
func NumerictoFloat(n pgtype.Numeric) float64 {
	if !n.Valid {
		log.Printf("Invalid numeric value")
		return 0.0
	}

	// Use Value() method to extract the value FROM pgtype.Numeric
	val, err := n.Value()
	if err != nil {
		log.Printf("failed to get numeric value: %v", err)
		return 0.0
	}

	// Value() returns a string for pgtype.Numeric
	if str, ok := val.(string); ok {
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			log.Printf("failed to parse numeric string %s to float: %v", str, err)
			return 0.0
		}
		log.Printf("Converted numeric to float: %v", f)
		return f
	}

	log.Printf("Unexpected value type from numeric: %T", val)
	return 0.0
}

// function to create the order with the listed items
func (r *Repo) CreateOrderWithListedItems(ctx context.Context, userid int, args []db.AddOrderItemParams) (db.Order, []db.OrderItem, error) {
	var createdorder db.Order
	var ordereditems []db.OrderItem
	// we initialize it with the transaction

	err := r.execTx(ctx, func(q *db.Queries) error {

		var total float64
		log.Printf("Starting order calculation with %d items", len(args))
		for _, arg := range args {
			itemPrice := NumerictoFloat(arg.Price)
			totalitem := itemPrice * float64(arg.Quantity)
			total += totalitem
			log.Printf("Item price: %f, quantity: %d, running total: %f",
				itemPrice, arg.Quantity, total)
		}

		totalamount := FloattoNumeric(total)
		//create the order
		ord, err := q.CreateOrder(ctx, db.CreateOrderParams{
			UserID:      int32(userid),
			TotalAmount: totalamount,
			StatusStaus: "pending",
		})
		if err != nil {
			log.Print("problem in creating order")
			return err
		}

		createdorder = ord
		//we need to insert the order items into the order itself that i just created

		for _, arg := range args {
			arg.OrderID = ord.OrderID
			ordereditem, err := q.AddOrderItem(ctx, arg)
			if err != nil {
				return err
			}
			ordereditems = append(ordereditems, ordereditem)
		}
		return nil

	})
	if err != nil {
		return db.Order{}, nil, err
	}

	//this is where i will inject it into the eventservice which is my service to handle the events

	err = r.eventservice.PublishOrderCreated(ctx, createdorder, ordereditems)
	if err != nil {
		log.Print("problem in publishing order created event")
		return db.Order{}, nil, err
	}

	return createdorder, ordereditems, nil

}

// func to list a users orders

func (r *Repo) ListOrders(ctx context.Context, args db.ListOrdersParams) ([]db.Order, error) {

	//listedItems, err := r.q.ListOrders(ctx, args)
	listedItems, err := r.q.ListOrders(ctx, args)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return listedItems, nil

}

// func to update the order status
func (r *Repo) UpdateOrderStatus(ctx context.Context, orderid int32, statusstatus string) (db.Order, error) {
	// Get the current order to check the old status
	currentOrder, err := r.GetOrder(ctx, orderid)
	if err != nil {
		log.Printf("error getting current order status: %v", err)
		return db.Order{}, err
	}
	oldstatusstatus := currentOrder.StatusStaus

	orderstatus, err := r.q.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		OrderID:     orderid,
		StatusStaus: statusstatus,
	})
	if err != nil {
		log.Printf("error updating order status for order %d: %v", orderid, err)
		return db.Order{}, err
	}

	if oldstatusstatus != orderstatus.StatusStaus {
		err = r.eventservice.PublishOrderStatusUpdated(ctx, orderstatus.OrderID, orderstatus.UserID, oldstatusstatus, orderstatus.StatusStaus)
		if err != nil {
			log.Printf("error publishing status update event for order %d: %v", orderid, err)
			return db.Order{}, err
		}
	}
	return orderstatus, nil
}

// func to get an order
func (r *Repo) GetOrder(ctx context.Context, orderid int32) (db.Order, error) {
	gottenorder, err := r.q.GetOrder(ctx, orderid)
	if err != nil {
		log.Printf("error getting order with ID %d: %v", orderid, err)
		return db.Order{}, err
	}
	return gottenorder, nil
}
