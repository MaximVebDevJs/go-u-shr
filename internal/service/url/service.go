package url

type service struct {
	urlRepo UrlRepository
	baseURL string
}

func New(urlRepo UrlRepository, baseUrl string) *service {
	return &service{urlRepo: urlRepo, baseURL: baseUrl}
}
