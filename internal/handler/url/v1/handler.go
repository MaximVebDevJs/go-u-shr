package v1

type Handler struct {
	urlService UrlService
}

func New(urlService UrlService) *Handler {
	return &Handler{urlService: urlService}
}
