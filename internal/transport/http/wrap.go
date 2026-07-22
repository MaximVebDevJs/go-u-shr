package http

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type ErrorHandlerFunc func(
	context.Context,
	*zap.Logger,
	http.ResponseWriter,
	*http.Request,
	error,
)

func Wrap(fn HandlerFunc, eh ErrorHandlerFunc, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			eh(r.Context(), logger, w, r, err)
		}
	}
}
