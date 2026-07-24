package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MaximVebDevJs/go-u-shr/internal/auth"
	serviceurl "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	transporthttp "github.com/MaximVebDevJs/go-u-shr/internal/transport/http"
	"go.uber.org/zap"
)

// batchRequestItem — DTO одного элемента входящего батча.
type batchRequestItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// batchResponseItem — DTO одного элемента ответа батча.
type batchResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// BatchUrls POST /api/shorten/batch.
func (h *Handler) BatchUrls(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, transporthttp.MaxRequestBodySize)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("не удалось прочитать тело запроса: %w", err)
	}

	var req []batchRequestItem
	if err = json.Unmarshal(body, &req); err != nil {
		return serviceurl.ErrInvalidJSON
	}

	items := make([]serviceurl.BatchItem, 0, len(req))
	for _, item := range req {
		correlationID := strings.TrimSpace(item.CorrelationID)
		originalURL := strings.TrimSpace(item.OriginalURL)

		if correlationID == "" || originalURL == "" {
			return serviceurl.ErrInvalidUrl
		}

		items = append(items, serviceurl.BatchItem{
			CorrelationID: correlationID,
			OriginalURL:   originalURL,
		})
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	results, err := h.urlService.BatchCreate(r.Context(), userID, items)
	if err != nil {
		var duplicateErr *serviceurl.DuplicateURLsError
		if errors.As(err, &duplicateErr) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)

			if encErr := json.NewEncoder(w).Encode(duplicateErr.URLs); encErr != nil {
				h.logger.Error("не удалось закодировать JSON-ответ", zap.Error(encErr))
			}

			return nil
		}

		return fmt.Errorf("ошибка пакетного создания коротких ссылок: %w", err)
	}

	h.logger.Info("пакет коротких ссылок успешно создан",
		zap.Int("batch_size", len(results)),
	)

	resp := make([]batchResponseItem, 0, len(results))
	for _, result := range results {
		resp = append(resp, batchResponseItem{
			CorrelationID: result.CorrelationID,
			ShortURL:      result.ShortURL,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if encErr := json.NewEncoder(w).Encode(resp); encErr != nil {
		// Заголовки уже отправлены — не возвращаем ошибку, чтобы Wrap не записал второй ответ.
		h.logger.Error("не удалось закодировать JSON-ответ", zap.Error(encErr))

		return nil
	}

	return nil
}
