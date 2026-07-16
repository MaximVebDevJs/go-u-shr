package url

import (
	"context"
)

type UrlRepository interface {
	Get(ctx context.Context, id string) (string, error)
	Create(ctx context.Context, url string) error
}
