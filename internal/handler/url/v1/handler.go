package v1

import "go.uber.org/zap"

type Handler struct {
	urlService UrlService
	logger     *zap.Logger
}

func New(urlService UrlService, logger *zap.Logger) *Handler {
	return &Handler{urlService: urlService, logger: logger}
}
