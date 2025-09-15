package internal

import (
	"context"

	db "github.com/franzego/store-service/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	dbq *db.Queries
	db  *pgxpool.Pool
}

func NewRepoService(dbconn *pgxpool.Pool) *Repo {
	return &Repo{
		dbq: db.New(dbconn),
		db:  dbconn,
	}
}

func (r *Repo) BulkCheckStock(ctx context.Context, arg db.BulkCheckStockParams)
