package url

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// upsertURLQuery вставляет alias либо возвращает уже занятый для этого original_url.
// Флаг owned = TRUE означает, что строка принадлежит текущему вызову: она либо только что
// вставлена, либо была помечена удалённой и восстановлена. Без восстановления UNIQUE
// (original_url) навсегда блокировал бы повторное сокращение удалённой ссылки.
const upsertURLQuery = `
	WITH upserted AS (
		INSERT INTO urls (original_url, uuid, user_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (original_url) DO UPDATE
			SET is_deleted = FALSE,
				user_id    = EXCLUDED.user_id
			WHERE urls.is_deleted
		RETURNING uuid, TRUE AS owned
	)
	SELECT uuid, owned
	FROM (
		SELECT uuid, owned FROM upserted
		UNION ALL
		SELECT uuid, FALSE FROM urls WHERE original_url = $1
	) AS candidates
	ORDER BY owned DESC
	LIMIT 1`

const selectIDByOriginalURL = `SELECT uuid FROM urls WHERE original_url = $1`

func (r *Repository) Create(ctx context.Context, userID string, originalURL string, id string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("создать запись: %w", err)
	}

	var (
		returnedID string
		owned      bool
	)

	err := r.pool.QueryRow(ctx, upsertURLQuery, originalURL, id, userID).Scan(&returnedID, &owned)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// ON CONFLICT DO UPDATE ждёт конкурирующую транзакцию, но её строка могла стать
			// видимой уже после снапшота нашего запроса — тогда SELECT ветка вернула 0 строк.
			return r.existingID(ctx, originalURL)
		}

		return "", mapCreateError(err)
	}

	if !owned {
		return returnedID, ErrOriginalURLExists
	}

	return returnedID, nil
}

// existingID возвращает alias, которым уже занят original_url, отдельным запросом:
// новый statement получает свежий снапшот и видит закоммиченную конкурентом строку.
func (r *Repository) existingID(ctx context.Context, originalURL string) (string, error) {
	var existingID string

	if err := r.pool.QueryRow(ctx, selectIDByOriginalURL, originalURL).Scan(&existingID); err != nil {
		return "", fmt.Errorf("создать url: получить занятый alias: %w", err)
	}

	return existingID, ErrOriginalURLExists
}
