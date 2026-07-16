package v1

import (
	"net/http"

	transporthttp "github.com/MaximVebDevJs/go-u-shr/internal/transport/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc(
		"POST /",
		transporthttp.Wrap(h.CreateUrl, transporthttp.ErrorHandler),
	)

	mux.HandleFunc(
		"GET /{id}",
		transporthttp.Wrap(h.GetUrl, transporthttp.ErrorHandler),
	)
}
