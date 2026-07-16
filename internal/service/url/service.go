package url

import "github.com/MaximVebDevJs/go-u-shr/internal/config"

type service struct {
	urlRepo UrlRepository
	baseURL string
}

func New(urlRepo UrlRepository) *service {
	return &service{urlRepo: urlRepo, baseURL: config.FlagHTTPAddr}
}
