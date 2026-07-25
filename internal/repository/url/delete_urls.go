package url

import (
	"context"
	"fmt"
)

// MarkDeleted помечает URL пользователя удалёнными одним batch update.
func (r *Repository) MarkDeleted(ctx context.Context, userID string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("пометить urls удалёнными: %w", err)
	}

	const query = `
		UPDATE urls
		SET is_deleted = TRUE
		WHERE user_id = $1
		  AND uuid = ANY($2::text[])
		  AND is_deleted = FALSE
	`

	// RowsAffected намеренно не проверяем: missing, чужие и уже удалённые id — идемпотентный no-op.
	if _, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, userID, ids); err != nil {
		return fmt.Errorf("пометить urls удалёнными: %w", err)
	}

	return nil
}
