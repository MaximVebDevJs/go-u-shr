package v1

import (
	"net/http"
)

// GetUrl GET /{id}.
func (h *Handler) GetUrl(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	originalURL, err := h.urlService.Get(r.Context(), id)
	if err != nil {
		return err
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

	return nil
}
