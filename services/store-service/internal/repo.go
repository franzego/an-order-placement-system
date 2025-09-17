package internal

import (
	"context"
	"log"

	db "github.com/franzego/store-service/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoServicer interface {
	CheckStock(ctx context.Context, id int64) (pgtype.Numeric, error)
	ReserveStock(ctx context.Context, args db.ReserveStockParams) error
	ReleaseStock(ctx context.Context, args db.ReleaseStockParams) error
	GetProduct(ctx context.Context, id int64) (db.GetProductRow, error)
}

type Repo struct {
	dbq *db.Queries
	db  *pgxpool.Pool
}

func NewRepoService(dbconn *pgxpool.Pool) *Repo {
	if dbconn == nil {
		return nil
	}
	return &Repo{
		dbq: db.New(dbconn),
		db:  dbconn,
	}
}

// we are returning the available_quantity
func (r *Repo) CheckStock(ctx context.Context, id int64) (pgtype.Numeric, error) {
	avaiQuant, err := r.dbq.CheckStock(ctx, id)
	if err != nil {
		log.Printf("couldn't run the function to check the stock in the database: %v", err)
	}
	return avaiQuant, err
}

func (r *Repo) ReserveStock(ctx context.Context, args db.ReserveStockParams) error {
	err := r.dbq.ReserveStock(ctx, args)
	if err != nil {
		log.Print(err)
	}
	return err
}

func (r *Repo) ReleaseStock(ctx context.Context, args db.ReleaseStockParams) error {
	err := r.dbq.ReleaseStock(ctx, args)
	if err != nil {
		log.Print(err)
	}
	return err
}

func (r *Repo) GetProduct(ctx context.Context, id int64) (db.GetProductRow, error) {
	product, err := r.dbq.GetProduct(ctx, id)
	if err != nil {
		log.Print(err)
	}
	return product, err
}
