package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/MaximVebDevJs/go-u-shr/internal/auth"
	"github.com/MaximVebDevJs/go-u-shr/internal/config"
	apiUrlV1 "github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1"
	urlRepo "github.com/MaximVebDevJs/go-u-shr/internal/repository/url"
	urlService "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	"github.com/MaximVebDevJs/go-u-shr/internal/transport/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// NewHTTPHandler собирает HTTP router и возвращает cleanup для фоновых задач приложения.
func NewHTTPHandler(cfg *config.Config, log *zap.Logger, pool *pgxpool.Pool) (http.Handler, func(context.Context) error, error) {
	if strings.TrimSpace(cfg.AuthSecret) == "" {
		return nil, nil, fmt.Errorf("auth secret is required")
	}

	repo := urlRepo.New(pool)

	svc := urlService.New(repo, cfg.BaseURL, log)
	handler := apiUrlV1.New(svc, log)
	signer := auth.NewSigner(cfg.AuthSecret)

	r := chi.NewRouter()
	// RequestID идёт первым: его значение попадает в лог запроса и в лог ошибки.
	r.Use(chiMiddleware.RequestID)
	// Recoverer перехватывает panic в handlers/middleware и отдаёт 500 вместо падения процесса.
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.DecompressMiddleware)
	r.Use(middleware.CompressMiddleware)
	r.Use(middleware.RequestLogger(log))

	apiUrlV1.RegisterRoutes(r, handler, log, signer)

	// Cleanup отдаём наружу, чтобы владелец HTTP-сервера остановил service-owned worker на shutdown.
	return r, svc.Close, nil
}
