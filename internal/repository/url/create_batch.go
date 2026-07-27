package url

import (
	"context"
	"errors"
	"fmt"

	"github.com/MaximVebDevJs/go-u-shr/internal/model"
	"github.com/jackc/pgx/v5"
)

// BatchRecord описывает одну запись для пакетной вставки.
type BatchRecord = model.BatchRecord

// CreateBatch сохраняет несколько URL пользователя в одной транзакции.
// Либо все строки записываются, либо ни одна (rollback при ошибке).
func (r *Repository) CreateBatch(ctx context.Context, userID string, records []BatchRecord) error {
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

	// Rollback нужен для отката при ошибке вставки или Commit.
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	conflicts, err := upsertBatch(ctx, tx, userID, records)
	if err != nil {
		return err
	}

	if len(conflicts) > 0 {
		return &model.DuplicateOriginalURLsError{URLs: conflicts}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("зафиксировать транзакцию: %w", err)
	}

	return nil
}

// upsertBatch отправляет все вставки одним batch-запросом и возвращает original_url,
// уже занятые живыми записями.
func upsertBatch(ctx context.Context, tx pgx.Tx, userID string, records []BatchRecord) ([]string, error) {
	batch := &pgx.Batch{}
	for _, rec := range records {
		batch.Queue(upsertURLQuery, rec.OriginalURL, rec.ID, userID)
	}

	results := tx.SendBatch(ctx, batch)

	var (
		conflicts = make([]string, 0)
		firstErr  error
	)

	for _, rec := range records {
		var (
			returnedID string
			owned      bool
		)

		err := results.QueryRow().Scan(&returnedID, &owned)

		switch {
		// ErrNoRows означает, что конкурирующая транзакция заняла original_url уже после
		// снапшота нашего запроса — для клиента это такой же конфликт.
		case errors.Is(err, pgx.ErrNoRows):
			conflicts = appendUniqueURL(conflicts, rec.OriginalURL)
		case err != nil:
			// Первая SQL-ошибка обрывает всю транзакцию, остальные результаты повторяют её.
			if firstErr == nil {
				firstErr = err
			}
		case !owned:
			conflicts = appendUniqueURL(conflicts, rec.OriginalURL)
		}
	}

	// Close обязателен до Commit: он дочитывает оставшиеся результаты batch.
	if closeErr := results.Close(); closeErr != nil && firstErr == nil {
		firstErr = closeErr
	}

	if firstErr != nil {
		return nil, mapCreateError(firstErr)
	}

	return conflicts, nil
}

func appendUniqueURL(urls []string, url string) []string {
	for _, existing := range urls {
		if existing == url {
			return urls
		}
	}

	return append(urls, url)
}
