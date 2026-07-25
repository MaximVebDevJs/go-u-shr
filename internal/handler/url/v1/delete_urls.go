package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MaximVebDevJs/go-u-shr/internal/auth"
	serviceurl "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	transporthttp "github.com/MaximVebDevJs/go-u-shr/internal/transport/http"
	"go.uber.org/zap"
)

// DeleteUrls DELETE /api/user/urls.
func (h *Handler) DeleteUrls(w http.ResponseWriter, r *http.Request) error {
	// Ограничиваем размер тела, чтобы один запрос не съел всю память процесса.
	r.Body = http.MaxBytesReader(w, r.Body, transporthttp.MaxRequestBodySize)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("не удалось прочитать тело запроса: %w", err)
	}

	var req []string
	if err = json.Unmarshal(body, &req); err != nil {
		return serviceurl.ErrInvalidJSON
	}
	if len(req) == 0 {
		return serviceurl.ErrInvalidJSON
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	if err = h.urlService.DeleteUrls(r.Context(), userID, req); err != nil {
		return fmt.Errorf("ошибка удаления коротких ссылок: %w", err)
	}

	// Логируем только агрегаты: сами short id — пользовательские данные и в логах не нужны.
	h.logger.Info("запрос на асинхронное удаление ссылок принят",
		zap.String("user_id", userID),
		zap.Int("input_count", len(req)),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	return nil
}
