package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zepyrshut/rating-orama/internal/app"
	"github.com/zepyrshut/rating-orama/internal/sqlc"
)

type pgxRepository struct {
	*sqlc.Queries
	pool *pgxpool.Pool
	app  *app.ExtendedApp
}

var _ ExtendedQuerier = &pgxRepository{}

func NewPGXRepo(pgx *pgxpool.Pool, app *app.ExtendedApp) ExtendedQuerier {
	return &pgxRepository{
		Queries: sqlc.New(pgx),
		pool:    pgx,
		app:     app,
	}
}

func (r *pgxRepository) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
