package v1

import (
	"github.com/MaximVebDevJs/go-u-shr/internal/auth"
	transporthttp "github.com/MaximVebDevJs/go-u-shr/internal/transport/http"
	"github.com/MaximVebDevJs/go-u-shr/internal/transport/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func RegisterRoutes(r chi.Router, h *Handler, log *zap.Logger, signer *auth.Signer) {
	r.With(middleware.EnsureAuthMiddleware(signer)).Post(
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
	r.With(middleware.MustAuth(signer)).Get(
		"/api/user/urls",
		transporthttp.Wrap(h.GetUserURLs, transporthttp.ErrorHandler, log),
	)
	// DELETE принимает задачу удаления только от аутентифицированного владельца cookie.
	r.With(middleware.MustAuth(signer)).Delete(
		"/api/user/urls",
		transporthttp.Wrap(h.DeleteUrls, transporthttp.ErrorHandler, log),
	)
	r.Route("/api/shorten", func(r chi.Router) {
		r.With(middleware.EnsureAuthMiddleware(signer)).Post(
			"/",
			transporthttp.Wrap(h.CreateUrlJSON, transporthttp.ErrorHandler, log))
		r.With(middleware.EnsureAuthMiddleware(signer)).Post(
			"/batch",
			transporthttp.Wrap(h.BatchUrls, transporthttp.ErrorHandler, log))
		r.Get(
			"/{id}",
			transporthttp.Wrap(h.GetUrl, transporthttp.ErrorHandler, log))
	})
}
