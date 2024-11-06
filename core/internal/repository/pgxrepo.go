package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zepyrshut/rating-orama/internal/sqlc"
)

type pgxRepository struct {
	*sqlc.Queries
	db *pgxpool.Pool
}

func NewPGXRepo(db *pgxpool.Pool) ExtendedQuerier {
	return &pgxRepository{
		Queries: sqlc.New(db),
		db:      db,
	}
}

func (r *pgxRepository) execTx(ctx context.Context, txFunc func(tx pgx.Tx) error) error {
	slog.Info("starting transaction", "txFunc", txFunc)
	tx, err := r.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start transaction", "error", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := txFunc(tx); err != nil {
		slog.Error("failed to execute transaction", "error", err)
		return err
	}

	slog.Info("committing transaction", "txFunc", txFunc)
	return tx.Commit(ctx)
}
