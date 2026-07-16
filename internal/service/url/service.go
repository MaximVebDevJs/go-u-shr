package url

type service struct {
	urlRepo UrlRepository
	baseURL string
}

func New(urlRepo UrlRepository) *service {
	return &service{urlRepo: urlRepo, baseURL: "0.0.0.0:8080"}
}
