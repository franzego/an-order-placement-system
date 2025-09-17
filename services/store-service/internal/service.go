package internal

import (
	"context"
	"fmt"
	"log"
	"strconv"

	db "github.com/franzego/store-service/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	RepoServicer
}

func NewService(svc RepoServicer) *Service {
	if svc == nil {
		return nil
	}
	return &Service{RepoServicer: svc}
}

//function to convert pgtypenumeric to float

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

// function to convert float back to numeric
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

// func to reseverve stock
func (s *Service) StoreFunc(ctx context.Context, Id int32, quantity int32) (db.GetProductRow, error) {
	availquant, err := s.RepoServicer.CheckStock(ctx, int64(Id))
	if err != nil {
		log.Print(err)
	}

	num := NumerictoFloat(availquant)
	/*numargs := NumerictoFloat(args.AvailableQuantity)
	if num < float64(quantity) {
		//log.Print("the available quantity is not enough")
		return db.GetProductRow{}, fmt.Errorf("not enough stock: have %v, need %v", num, numargs)

	}*/
	pgavail := FloattoNumeric(num)

	sth := db.ReserveStockParams{
		AvailableQuantity: pgavail,
		ID:                int64(Id),
	}
	err = s.RepoServicer.ReserveStock(ctx, sth)
	if err != nil {
		log.Print("There was an error in reserving the item")
		return db.GetProductRow{}, fmt.Errorf("failed to reserve stock: %w", err)
	}

	/*err = s.RepoServicer.ReleaseStock(ctx, db.ReleaseStockParams(sth))
	if err != nil {
		log.Print("There was an error in releasing the item")
	}*/
	product, err := s.RepoServicer.GetProduct(ctx, sth.ID)
	if err != nil {
		log.Print("There was an error in getting the item")
		return db.GetProductRow{}, fmt.Errorf("failed to get product after reserve: %w", err)
	}
	return product, nil

}
