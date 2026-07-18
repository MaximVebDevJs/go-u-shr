package app

import (
	"net/http"

	"github.com/MaximVebDevJs/go-u-shr/internal/config"
	apiUrlV1 "github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1"
	"github.com/MaximVebDevJs/go-u-shr/internal/logger"
	urlRepo "github.com/MaximVebDevJs/go-u-shr/internal/repository/url"
	urlService "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	"github.com/MaximVebDevJs/go-u-shr/internal/transport/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewHTTPHandler(cfg *config.Config, log *zap.Logger) http.Handler {
	repo := urlRepo.New(cfg.FileStoragePath)
	svc := urlService.New(repo, cfg.BaseURL)
	handler := apiUrlV1.New(svc, log)

	r := chi.NewRouter()
	r.Use(middleware.DecompressMiddleware)
	r.Use(middleware.CompressMiddleware)
	r.Use(logger.RequestLogger(log))

	apiUrlV1.RegisterRoutes(r, handler)

	return r
}
