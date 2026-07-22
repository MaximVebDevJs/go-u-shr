package v1

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	serviceurl "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	transporthttp "github.com/MaximVebDevJs/go-u-shr/internal/transport/http"
	"go.uber.org/zap"
)

// CreateUrl POST /.
func (h *Handler) CreateUrl(w http.ResponseWriter, r *http.Request) error {
	// Ограничиваем размер тела, чтобы один запрос не съел всю память процесса.
	r.Body = http.MaxBytesReader(w, r.Body, transporthttp.MaxRequestBodySize)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("не удалось прочитать тело запроса: %w", err)
	}

	originalURL := strings.TrimSpace(string(body))

	shortURL, err := h.urlService.Create(r.Context(), originalURL)
	if err != nil {
		if errors.Is(err, serviceurl.ErrURLAlreadyExists) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)

			if _, writeErr := w.Write([]byte(shortURL)); writeErr != nil {
				h.logger.Error("не удалось записать тело ответа", zap.Error(writeErr))
			}

			return nil
		}

		return fmt.Errorf("ошибка создания короткой ссылки: %w", err)
	}

	h.logger.Info("короткая ссылка успешно создана",
		append(urlLogFields(originalURL), zap.String("short", shortURL))...,
	)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	if _, writeErr := w.Write([]byte(shortURL)); writeErr != nil {
		// Заголовки уже отправлены — не возвращаем ошибку, чтобы Wrap не записал второй ответ.
		h.logger.Error("не удалось записать тело ответа", zap.Error(writeErr))

		return nil
	}

	return nil
}
