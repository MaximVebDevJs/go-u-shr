package v1

import (
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// CreateUrl POST /.
func (h *Handler) CreateUrl(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("не удалось прочитать тело запроса: %w", err)
	}

	shortURL, err := h.urlService.Create(r.Context(), string(body))
	if err != nil {
		return fmt.Errorf("ошибка создания короткой ссылки: %w", err)
	}

	h.logger.Info("короткая ссылка успешно создана",
		zap.String("original", string(body)),
		zap.String("short", shortURL),
	)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(shortURL))
	return err
}
