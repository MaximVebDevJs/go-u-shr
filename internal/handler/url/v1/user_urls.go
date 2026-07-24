package v1

import (
	"encoding/json"
	"net/http"

	"github.com/MaximVebDevJs/go-u-shr/internal/auth"
	"go.uber.org/zap"
)

type userURLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (h *Handler) GetUserURLs(w http.ResponseWriter, r *http.Request) error {
	userID, ok := auth.UserIDFromContext(
		r.Context(),
	)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	urls, err := h.urlService.GetUserURLs(
		r.Context(),
		userID,
	)

	if err != nil {
		return err
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	resp := make([]userURLResponse, 0, len(urls))
	for _, item := range urls {
		resp = append(resp, userURLResponse{
			ShortURL:    item.ShortURL,
			OriginalURL: item.OriginalURL,
		})
	}

	h.logger.Info("ссылки пользователя получены",
		zap.Int("batch_size", len(urls)),
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	return json.NewEncoder(w).Encode(resp)
}
