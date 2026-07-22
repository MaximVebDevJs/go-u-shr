package v1

import (
	transporthttp "github.com/MaximVebDevJs/go-u-shr/internal/transport/http"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func RegisterRoutes(r chi.Router, h *Handler, log *zap.Logger) {
	r.Post(
		"/",
		transporthttp.Wrap(h.CreateUrl, transporthttp.ErrorHandler, log),
	)
	r.Get(
		"/ping",
		transporthttp.Wrap(h.Ping, transporthttp.ErrorHandler, log),
	)
	r.Get(
		"/{id}",
		transporthttp.Wrap(h.GetUrl, transporthttp.ErrorHandler, log),
	)
	r.Route("/api/shorten", func(r chi.Router) {
		r.Post(
			"/",
			transporthttp.Wrap(h.CreateUrlJSON, transporthttp.ErrorHandler, log))
		r.Get(
			"/{id}",
			transporthttp.Wrap(h.GetUrl, transporthttp.ErrorHandler, log))
	})
}
