package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MaximVebDevJs/go-u-shr/internal/auth"
	serviceurl "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	transporthttp "github.com/MaximVebDevJs/go-u-shr/internal/transport/http"
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
	// Ограничиваем размер тела, чтобы один запрос не съел всю память процесса.
	r.Body = http.MaxBytesReader(w, r.Body, transporthttp.MaxRequestBodySize)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("не удалось прочитать тело запроса: %w", err)
	}

	var req shortenRequest
	if err = json.Unmarshal(body, &req); err != nil {
		return serviceurl.ErrInvalidJSON
	}

	// TrimSpace здесь и в validateURL — защита от whitespace-only значений.
	req.URL = strings.TrimSpace(req.URL)
	if req.URL == "" {
		return serviceurl.ErrInvalidUrl
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	shortURL, err := h.urlService.Create(r.Context(), userID, req.URL)
	if err != nil {
		if errors.Is(err, serviceurl.ErrURLAlreadyExists) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)

			if encErr := json.NewEncoder(w).Encode(shortenResponse{Result: shortURL}); encErr != nil {
				h.logger.Error("не удалось закодировать JSON-ответ", zap.Error(encErr))
			}

			return nil
		}

		return fmt.Errorf("ошибка создания короткой ссылки: %w", err)
	}

	h.logger.Info("короткая ссылка успешно создана",
		append(urlLogFields(req.URL), zap.String("short", shortURL))...,
	)

	resp := shortenResponse{Result: shortURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if encErr := json.NewEncoder(w).Encode(resp); encErr != nil {
		// Заголовки уже отправлены — не возвращаем ошибку, чтобы Wrap не записал второй ответ.
		h.logger.Error("не удалось закодировать JSON-ответ", zap.Error(encErr))

		return nil
	}

	return nil
}
