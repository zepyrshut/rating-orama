package db

import (
	"context"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func NewDBPool(dataSource string) *pgxpool.Pool {
	dbPool, err := pgxpool.New(context.Background(), dataSource)
	if err != nil {
		log.Fatal(err)
	}

	return dbPool
}
