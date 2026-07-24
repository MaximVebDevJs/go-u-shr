package url

import (
	"context"
	"fmt"

	"github.com/MaximVebDevJs/go-u-shr/internal/model"
)

// GetUserURLs возвращает все URL, созданные пользователем.
func (r *Repository) GetUserURLs(
	ctx context.Context,
	userID string,
) ([]model.UserURL, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("получить ссылки пользователя: %w", err)
	}

	rows, err := r.pool.Query(
		ctx,
		`
		SELECT uuid, original_url
		FROM urls
		WHERE user_id=$1
		ORDER BY id
		`,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("получить ссылки пользователя: %w", err)
	}

	defer rows.Close()

	result := make([]model.UserURL, 0)

	for rows.Next() {

		var item model.UserURL

		err := rows.Scan(
			&item.ShortID,
			&item.OriginalURL,
		)

		if err != nil {
			return nil, fmt.Errorf("прочитать ссылки пользователя: %w", err)
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("прочитать ссылки пользователя: %w", err)
	}

	return result, nil
}
