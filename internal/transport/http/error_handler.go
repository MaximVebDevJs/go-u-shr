package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	svc "github.com/MaximVebDevJs/go-u-shr/internal/service/url"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type errorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func ErrorHandler(ctx context.Context, logger *zap.Logger, w http.ResponseWriter, r *http.Request, err error) {
	code, message := mapError(err)

	logger.Error(
		"ошибка запроса",
		zap.Error(err),
		zap.Int("status", code),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("request_id", chiMiddleware.GetReqID(ctx)),
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
			zap.String("request_id", chiMiddleware.GetReqID(ctx)),
		)
	}
}

// mapError возвращает публичный статус и сообщение. Наружу отдаём текст только известных
// доменных ошибок: обёрнутая цепочка раскрыла бы клиенту внутреннее устройство слоёв.
func mapError(err error) (int, string) {
	switch {
	// 404 Not Found
	case errors.Is(err, svc.ErrUrlNotFound):
		return http.StatusNotFound, svc.ErrUrlNotFound.Error()

	// 410 Gone
	case errors.Is(err, svc.ErrURLDeleted):
		return http.StatusGone, svc.ErrURLDeleted.Error()

	// 400 Bad Request
	case errors.Is(err, svc.ErrInvalidUrl):
		return http.StatusBadRequest, svc.ErrInvalidUrl.Error()
	case errors.Is(err, svc.ErrInvalidJSON):
		return http.StatusBadRequest, svc.ErrInvalidJSON.Error()

	// 413 Payload Too Large
	case errors.Is(err, svc.ErrBatchTooLarge):
		return http.StatusRequestEntityTooLarge, svc.ErrBatchTooLarge.Error()

	// 429 Too Many Requests — сервис не успевает разгружать очередь удаления.
	case errors.Is(err, svc.ErrDeleteQueueFull):
		return http.StatusTooManyRequests, svc.ErrDeleteQueueFull.Error()
	}

	// 413 Payload Too Large — срабатывает при превышении MaxBytesReader.
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		return http.StatusRequestEntityTooLarge, "слишком большое тело запроса"
	}

	// 500 Internal Server Error
	return http.StatusInternalServerError, "внутренняя ошибка"
}
