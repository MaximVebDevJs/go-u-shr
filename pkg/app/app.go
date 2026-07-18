package app

import (
	"net/http"

	apiUrlV1 "github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1"
	urlRepo "github.com/MaximVebDevJs/go-u-shr/internal/repository/url"
	urlService "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHTTPHandler(baseUrl string) http.Handler {
	repo := urlRepo.New()
	svc := urlService.New(repo, baseUrl)
	handler := apiUrlV1.New(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	apiUrlV1.RegisterRoutes(r, handler)

	return r
}
