package v1

import (
	"encoding/json"
	"io"
	"net/http"

	errs "github.com/MaximVebDevJs/go-u-shr/internal/errors"
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
		h.logger.Error("не удалось прочитать тело запроса", zap.Error(err))
		return err
	}

	var req shortenRequest
	if err = json.Unmarshal(body, &req); err != nil {
		h.logger.Error("передан невалидный JSON", zap.Error(err))
		return errs.ErrInvalidJSON
	}

	if req.URL == "" {
		h.logger.Error("передана пустая ссылка", zap.String("url", req.URL))
		return errs.ErrInvalidUrl
	}

	shortURL, err := h.urlService.Create(r.Context(), req.URL)
	if err != nil {
		h.logger.Error("ошибка создания короткой ссылки",
			zap.Error(err),
			zap.String("original_url", string(body)),
		)
		return err
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
