package url

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type queryRower interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func (r *Repository) Create(ctx context.Context, userID string, originalURL string, id string) (string, error) {
	return r.insertURL(ctx, r.getter.DefaultTrOrDB(ctx, r.pool), userID, originalURL, id)
}

func (r *Repository) insertURL(
	ctx context.Context,
	q queryRower,
	userID string,
	originalURL string,
	id string,
) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("создать запись: %w", err)
	}

	const query = `
		WITH inserted AS (
			INSERT INTO urls (original_url, uuid, user_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (original_url) DO NOTHING
			RETURNING uuid
		)
		SELECT uuid FROM inserted
		UNION ALL
		SELECT uuid FROM urls WHERE original_url = $1
		LIMIT 1`

	var returnedID string

	err := q.QueryRow(ctx, query, originalURL, id, userID).Scan(&returnedID)
	if err != nil {
		return "", mapCreateError(err)
	}

	if returnedID != id {
		return returnedID, ErrOriginalURLExists
	}

	return returnedID, nil
}
