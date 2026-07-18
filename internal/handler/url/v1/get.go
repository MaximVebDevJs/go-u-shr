package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// GetUrl GET /{id}.
func (h *Handler) GetUrl(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	originalURL, err := h.urlService.Get(r.Context(), id)
	if err != nil {
		h.logger.Error("ошибка получения короткой ссылки",
			zap.Error(err),
			zap.String("link_id", id),
		)

		return err
	}

	h.logger.Info("короткая ссылка успешно получена",
		zap.String("id", id),
		zap.String("original", originalURL),
	)

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

	return nil
}
