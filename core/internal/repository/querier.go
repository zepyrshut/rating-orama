package repository

import (
	"github.com/zepyrshut/rating-orama/internal/sqlc"
)

type ExtendedQuerier interface {
	sqlc.Querier
}
