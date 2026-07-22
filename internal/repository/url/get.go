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
		SELECT original_url FROM urls WHERE uuid = $1
	`

	var originalURL string

	err := r.getter.DefaultTrOrDB(ctx, r.pool).QueryRow(ctx, query, id).Scan(&originalURL)
	if err != nil {
		return "", mapGetError(err)
	}

	return originalURL, nil
}
