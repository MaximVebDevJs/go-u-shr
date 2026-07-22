package url

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}
