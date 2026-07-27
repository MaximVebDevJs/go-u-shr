package url

import (
	"context"
	"fmt"
)

func (r *Repository) Get(ctx context.Context, id string) (string, error) {
	// (например, клиент закрыл соединение или вышел таймаут).
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("получить запись: %w", err)
	}

	const query = `
		SELECT original_url, is_deleted FROM urls WHERE uuid = $1
	`

	var originalURL string
	var isDeleted bool

	err := r.pool.QueryRow(ctx, query, id).Scan(&originalURL, &isDeleted)
	if err != nil {
		return "", mapGetError(err)
	}

	if isDeleted {
		// Удалённая запись существует в БД, но публичный GET должен отвечать 410 Gone.
		return "", ErrDeleted
	}

	return originalURL, nil
}
