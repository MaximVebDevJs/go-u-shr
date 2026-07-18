package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	errs "github.com/MaximVebDevJs/go-u-shr/internal/errors"
)

type errorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func ErrorHandler(ctx context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	code, message := mapError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if encErr := json.NewEncoder(w).Encode(errorResponse{
		Code:    code,
		Message: message,
	}); encErr != nil {
		slog.ErrorContext(ctx, "ошибка кодирования ответа", "error", encErr)
	}
}

func mapError(err error) (int, string) {
	switch {
	// 404 Not Found
	case errors.Is(err, errs.ErrUrlNotFound):
		return http.StatusNotFound, err.Error()

	// 400 Bad Request
	case errors.Is(err, errs.ErrInvalidUrl):
		return http.StatusBadRequest, err.Error()

	// 500 Internal Server Error
	default:
		return http.StatusInternalServerError, "внутренняя ошибка"
	}
}
