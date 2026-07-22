package url

import (
	"context"
	"fmt"
)

func (r *Repository) Create(ctx context.Context, originalUrl string, id string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("создать запись: %w", err)
	}

	const query = `
		INSERT INTO urls (original_url, uuid)
		VALUES ($1, $2)`

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(
		ctx,
		query,
		originalUrl,
		id,
	)
	if err != nil {
		return mapCreateError(err)
	}

	return nil
}
