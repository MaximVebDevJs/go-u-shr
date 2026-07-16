package http

import (
	"context"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type ErrorHandlerFunc func(
	context.Context,
	http.ResponseWriter,
	*http.Request,
	error,
)

func Wrap(fn HandlerFunc, eh ErrorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			eh(r.Context(), w, r, err)
		}
	}
}
