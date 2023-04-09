package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresDBRepo struct {
	DB *pgxpool.Pool
}

func NewPostgresRepo(conn *pgxpool.Pool) DBRepo {
	return &postgresDBRepo{
		DB: conn,
	}
}
