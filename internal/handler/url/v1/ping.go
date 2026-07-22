package v1

import (
	"net/http"

	"go.uber.org/zap"
)

// Ping GET /ping.
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) error {
	if err := h.urlService.Ping(r.Context()); err != nil {
		h.logger.Error("проверка базы данных не удалась",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Error(err),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return nil
	}

	h.logger.Info("проверка базы данных успешна",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	w.WriteHeader(http.StatusOK)

	return nil
}
