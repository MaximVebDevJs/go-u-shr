package url

import (
	"context"
	"fmt"
	"strings"
)

// BatchRecord описывает одну запись для пакетной вставки.
type BatchRecord struct {
	OriginalURL string
	ID          string
}

// CreateBatch сохраняет несколько URL в одной транзакции.
// Либо все строки записываются, либо ни одна (rollback при ошибке).
func (r *Repository) CreateBatch(ctx context.Context, records []BatchRecord) error {
	if len(records) == 0 {
		return nil
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("создать batch: %w", err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("начать транзакцию: %w", err)
	}

	// Rollback после успешного Commit — no-op; нужен для отката при ошибке Exec/Commit.
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// собираем sql запрос
	query, args := buildBatchInsertQuery(records)

	// выполняем sql запрос
	if _, err = tx.Exec(ctx, query, args...); err != nil {
		// маппим ошибку
		return mapCreateError(err)
	}

	// коммитим транзакцию
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("зафиксировать транзакцию: %w", err)
	}

	return nil
}

func buildBatchInsertQuery(records []BatchRecord) (string, []any) {
	// для конкатенации без цикла
	var sb strings.Builder
	// Да, по смыслу — это «допиши эту строку в конец», но не буквально sb +=
	sb.WriteString("INSERT INTO urls (original_url, uuid) VALUES ")

	args := make([]any, 0, len(records)*2)

	for i, rec := range records {
		if i > 0 {
			sb.WriteString(", ")
		}

		// дописываем аргументы для запроса
		sb.WriteString(fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		args = append(args, rec.OriginalURL, rec.ID)
	}

	return sb.String(), args
}
