package url

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const createTableSQL = `
CREATE TABLE IF NOT EXISTS urls (
	id           BIGSERIAL PRIMARY KEY,
	original_url TEXT NOT NULL UNIQUE,
	uuid         VARCHAR(255) NOT NULL UNIQUE
);`

const createOriginalURLUniqueIndexSQL = `
CREATE UNIQUE INDEX IF NOT EXISTS urls_original_url_uidx ON urls (original_url);`

// Migrate создаёт таблицу urls, если её ещё нет.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("миграция urls: %w", err)
	}

	if _, err := pool.Exec(ctx, createTableSQL); err != nil {
		return fmt.Errorf("миграция urls: %w", err)
	}

	if _, err := pool.Exec(ctx, createOriginalURLUniqueIndexSQL); err != nil {
		return fmt.Errorf("миграция urls: %w", err)
	}

	return nil
}
