package url

import (
	"context"
	"errors"
	"fmt"
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

	// Rollback нужен для отката при ошибке Exec/Commit.
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// создаем слайс для хранения дублирующихся original_url
	conflicts := make([]string, 0)

	for _, rec := range records {
		// вставляем запись в БД
		returnedID, insertErr := r.insertURL(ctx, tx, rec.OriginalURL, rec.ID)
		if errors.Is(insertErr, ErrOriginalURLExists) {
			conflicts = appendUniqueURL(conflicts, rec.OriginalURL)
			continue
		}

		if insertErr != nil {
			return insertErr
		}

		// если returnedID не равен rec.ID, то запись уже существует
		if returnedID != rec.ID {
			// добавляем original_url в слайс conflicts
			conflicts = appendUniqueURL(conflicts, rec.OriginalURL)
		}
	}

	if len(conflicts) > 0 {
		return &DuplicateOriginalURLsError{URLs: conflicts}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("зафиксировать транзакцию: %w", err)
	}

	return nil
}

func appendUniqueURL(urls []string, url string) []string {
	for _, existing := range urls {
		if existing == url {
			return urls
		}
	}

	return append(urls, url)
}
