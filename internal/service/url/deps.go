package url

import (
	"context"

	"github.com/MaximVebDevJs/go-u-shr/internal/model"
)

type UrlRepository interface {
	Get(ctx context.Context, url string) (model.Url, error)
	Create(ctx context.Context, url model.Url) (model.Url, error)
}
