package url

import (
	"context"
	"errors"
	"fmt"

	urlRepo "github.com/MaximVebDevJs/go-u-shr/internal/repository/url"
)

const maxBatchSize = 100

// BatchItem описывает один URL для пакетного сокращения.
type BatchItem struct {
	CorrelationID string
	OriginalURL   string
}

// BatchResult описывает результат сокращения одного URL из батча.
type BatchResult struct {
	CorrelationID string
	ShortURL      string
}

// BatchCreate сокращает несколько URL атомарно: pre-validation, затем одна транзакция в БД.
func (s *service) BatchCreate(ctx context.Context, items []BatchItem) ([]BatchResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("создать batch: %w", err)
	}

	if len(items) == 0 {
		return []BatchResult{}, nil
	}

	if len(items) > maxBatchSize {
		return nil, ErrInvalidUrl
	}

	// При первой ошибке — отклоняем весь запрос (partial success не поддерживается контрактом).
	for _, item := range items {
		if err := validateURL(item.OriginalURL); err != nil {
			return nil, err
		}
	}

	// При коллизии uuid в БД — перегенерируем весь батч и повторяем (как в Create).
	if duplicateURLs := findInBatchDuplicateURLs(items); len(duplicateURLs) > 0 {
		return nil, &DuplicateURLsError{URLs: duplicateURLs}
	}

	for range maxIDGenerationAttempts {
		// собираем записи и результаты
		records, results, err := s.buildBatchRecords(items)
		if err != nil {
			return nil, err
		}

		// сохраняем записи в БД
		err = s.urlRepo.CreateBatch(ctx, records)
		if err == nil {
			return results, nil
		}

		var duplicateRepo *urlRepo.DuplicateOriginalURLsError
		if errors.As(err, &duplicateRepo) {
			return nil, &DuplicateURLsError{URLs: duplicateRepo.URLs}
		}

		if errors.Is(err, urlRepo.ErrAlreadyExists) {
			continue
		}

		return nil, fmt.Errorf("сохранение batch: %w", err)
	}

	return nil, fmt.Errorf(
		"создать batch: не удалось сгенерировать уникальные id за %d попыток",
		maxIDGenerationAttempts,
	)
}

func findInBatchDuplicateURLs(items []BatchItem) []string {
	counts := make(map[string]int, len(items))
	for _, item := range items {
		counts[item.OriginalURL]++
	}

	duplicates := make([]string, 0)
	for originalURL, count := range counts {
		if count > 1 {
			duplicates = append(duplicates, originalURL)
		}
	}

	return duplicates
}

func (s *service) buildBatchRecords(items []BatchItem) ([]urlRepo.BatchRecord, []BatchResult, error) {

	// создаем map для хранения использованных shortUrl
	usedIDs := make(map[string]struct{}, len(items))
	// создаем слайс для записей
	records := make([]urlRepo.BatchRecord, 0, len(items))
	// создаем слайс для результатов
	results := make([]BatchResult, 0, len(items))

	for _, item := range items {
		// генерируем shortUrl
		id, err := generateUniqueBatchID(usedIDs)
		if err != nil {
			return nil, nil, err
		}

		records = append(records, urlRepo.BatchRecord{
			OriginalURL: item.OriginalURL,
			ID:          id,
		})
		results = append(results, BatchResult{
			CorrelationID: item.CorrelationID,
			ShortURL:      s.baseURL + "/" + id,
		})
	}

	return records, results, nil
}

// generateUniqueBatchID возвращает alias, уникальный внутри текущего батча.
func generateUniqueBatchID(usedIDs map[string]struct{}) (string, error) {
	for range maxIDGenerationAttempts {
		id := generateShortID()
		if _, exists := usedIDs[id]; exists {
			continue
		}

		usedIDs[id] = struct{}{}

		return id, nil
	}

	return "", fmt.Errorf("не удалось сгенерировать уникальный id внутри batch")
}
