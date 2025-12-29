package repository

import (
	"context"

	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/dao"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/database"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	dao.Querier
	StartTransaction(context.Context) (pgx.Tx, error)
	WithTx(tx pgx.Tx) dao.Querier
}

type repository struct {
	*dao.Queries
	db database.Database
}

func New(db database.Database) (Repository, error) {
	tx, err := db.Tx()
	if err != nil {
		return nil, err
	}
	return &repository{
		Queries: dao.New(tx),
		db:      db,
	}, nil
}

func (r *repository) StartTransaction(ctx context.Context) (pgx.Tx, error) {
	pool, err := r.db.Tx()
	if err != nil {
		return nil, err
	}
	return pool.Begin(ctx)
}

func (r *repository) WithTx(tx pgx.Tx) dao.Querier {
	return dao.New(tx)
}
