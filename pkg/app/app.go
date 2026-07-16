package app

import (
	"net/http"

	apiUrlV1 "github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1"
	urlRepo "github.com/MaximVebDevJs/go-u-shr/internal/repository/url"
	urlService "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
)

func NewHTTPHandler() http.Handler {
	repo := urlRepo.New()
	svc := urlService.New(repo)
	handler := apiUrlV1.New(svc)

	mux := http.NewServeMux()
	apiUrlV1.RegisterRoutes(mux, handler)

	return mux
}
