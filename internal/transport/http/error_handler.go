package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	svc "github.com/MaximVebDevJs/go-u-shr/internal/service/url"

	"go.uber.org/zap"
)

type errorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func ErrorHandler(ctx context.Context, logger *zap.Logger, w http.ResponseWriter, _ *http.Request, err error) {
	code, message := mapError(err)

	logger.Error(
		"ошибка запроса",
		zap.Error(err),
		zap.Int("status", code),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if encErr := json.NewEncoder(w).Encode(errorResponse{
		Code:    code,
		Message: message,
	}); encErr != nil {
		logger.Error(
			"ошибка кодирования ответа",
			zap.Error(encErr),
		)
	}
}

func mapError(err error) (int, string) {
	switch {
	// 404 Not Found
	case errors.Is(err, svc.ErrUrlNotFound):
		return http.StatusNotFound, err.Error()

	// 410 Gone
	case errors.Is(err, svc.ErrURLDeleted):
		return http.StatusGone, err.Error()

	// 400 Bad Request
	case errors.Is(err, svc.ErrInvalidUrl):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, svc.ErrInvalidJSON):
		return http.StatusBadRequest, err.Error()
	}

	// 413 Payload Too Large — срабатывает при превышении MaxBytesReader.
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		return http.StatusRequestEntityTooLarge, "слишком большое тело запроса"
	}

	// 500 Internal Server Error
	return http.StatusInternalServerError, "внутренняя ошибка"
}
