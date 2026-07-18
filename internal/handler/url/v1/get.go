package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetUrl GET /{id}.
func (h *Handler) GetUrl(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	originalURL, err := h.urlService.Get(r.Context(), id)
	if err != nil {
		return err
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

	return nil
}
