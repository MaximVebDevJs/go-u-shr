package v1

import (
	"io"
	"net/http"

	"go.uber.org/zap"
)

// CreateUrl POST /.
func (h *Handler) CreateUrl(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("не удалось прочитать тело запроса", zap.Error(err))
		return err
	}

	shortURL, err := h.urlService.Create(r.Context(), string(body))
	if err != nil {
		h.logger.Error("ошибка создания короткой ссылки",
			zap.Error(err),
			zap.String("original_url", string(body)),
		)
		return err
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
