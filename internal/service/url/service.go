package url

type service struct {
	urlRepo UrlRepository
}

func New(urlRepo UrlRepository) *service {
	return &service{urlRepo: urlRepo}
}
