package url

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type queryRower interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func (r *Repository) Create(ctx context.Context, originalURL string, id string) (string, error) {
	return r.insertURL(ctx, r.getter.DefaultTrOrDB(ctx, r.pool), originalURL, id)
}

func (r *Repository) insertURL(
	ctx context.Context,
	q queryRower,
	originalURL string,
	id string,
) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("создать запись: %w", err)
	}

	const query = `
		INSERT INTO urls (original_url, uuid)
		VALUES ($1, $2)
		ON CONFLICT (original_url) DO UPDATE SET original_url = urls.original_url
		RETURNING uuid`

	var returnedID string

	err := q.QueryRow(ctx, query, originalURL, id).Scan(&returnedID)
	if err != nil {
		return "", mapCreateError(err)
	}

	if returnedID != id {
		return returnedID, ErrOriginalURLExists
	}

	return returnedID, nil
}
