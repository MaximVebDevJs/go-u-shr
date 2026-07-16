package v1

import (
	"io"
	"net/http"
)

// CreateUrl POST /.
func (h *Handler) CreateUrl(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	shortURL, err := h.urlService.Create(r.Context(), string(body))
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(shortURL))
	return err
}
