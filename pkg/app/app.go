package app

import (
	"net/http"

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

func NewHTTPHandler(cfg *config.Config, log *zap.Logger, pool *pgxpool.Pool) (http.Handler, error) {
	repo := urlRepo.New(pool)

	svc := urlService.New(repo, cfg.BaseURL)
	handler := apiUrlV1.New(svc, log)

	r := chi.NewRouter()
	// Recoverer перехватывает panic в handlers/middleware и отдаёт 500 вместо падения процесса.
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.DecompressMiddleware)
	r.Use(middleware.CompressMiddleware)
	r.Use(middleware.RequestLogger(log))

	apiUrlV1.RegisterRoutes(r, handler, log)

	return r, nil
}
