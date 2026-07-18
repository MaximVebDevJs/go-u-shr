package v1

import (
	transporthttp "github.com/MaximVebDevJs/go-u-shr/internal/transport/http"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Post(
		"/",
		transporthttp.Wrap(h.CreateUrl, transporthttp.ErrorHandler),
	)

	r.Get(
		"/{id}",
		transporthttp.Wrap(h.GetUrl, transporthttp.ErrorHandler),
	)
}
