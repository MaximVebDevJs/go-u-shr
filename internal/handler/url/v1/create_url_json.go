package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	serviceurl "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	"go.uber.org/zap"
)

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Result string `json:"result"`
}

// CreateUrlJSON POST /api/shorten
func (h *Handler) CreateUrlJSON(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("не удалось прочитать тело запроса: %w", err)
	}

	var req shortenRequest
	if err = json.Unmarshal(body, &req); err != nil {
		return serviceurl.ErrInvalidJSON
	}

	if req.URL == "" {
		return serviceurl.ErrInvalidUrl
	}

	shortURL, err := h.urlService.Create(r.Context(), req.URL)
	if err != nil {
		return fmt.Errorf("ошибка создания короткой ссылки: %w", err)
	}

	h.logger.Info("короткая ссылка успешно создана",
		zap.String("original", string(body)),
		zap.String("short", shortURL),
	)

	resp := shortenResponse{Result: shortURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(resp)
}
