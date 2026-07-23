package v1

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// Ping GET /ping.
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) error {
	if err := h.urlService.Ping(r.Context()); err != nil {
		return fmt.Errorf("проверка базы данных: %w", err)
	}

	h.logger.Info("проверка базы данных успешна",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	w.WriteHeader(http.StatusOK)

	return nil
}
