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
	uuid         VARCHAR(255) NOT NULL UNIQUE,
	user_id      TEXT NOT NULL DEFAULT '',
	is_deleted   BOOLEAN NOT NULL DEFAULT FALSE
);`

const addUserIDColumnSQL = `
ALTER TABLE urls
ADD COLUMN IF NOT EXISTS user_id TEXT NOT NULL DEFAULT '';`

const addIsDeletedColumnSQL = `
ALTER TABLE urls
ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN NOT NULL DEFAULT FALSE;`

const createOriginalURLUniqueIndexSQL = `
CREATE UNIQUE INDEX IF NOT EXISTS urls_original_url_uidx ON urls (original_url);`

// Частичный индекс покрывает и выборку ссылок пользователя, и batch update при удалении:
// оба запроса фильтруют по user_id и is_deleted = FALSE.
const createUserIDIndexSQL = `
CREATE INDEX IF NOT EXISTS urls_user_id_idx ON urls (user_id) WHERE is_deleted = FALSE;`

// Migrate создаёт таблицу urls, если её ещё нет.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("миграция urls: %w", err)
	}

	statements := []string{
		createTableSQL,
		addUserIDColumnSQL,
		addIsDeletedColumnSQL,
		createOriginalURLUniqueIndexSQL,
		createUserIDIndexSQL,
	}

	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement); err != nil {
			return fmt.Errorf("миграция urls: %w", err)
		}
	}

	return nil
}
