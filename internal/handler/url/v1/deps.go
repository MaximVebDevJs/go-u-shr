package v1

import (
	"context"
)

type UrlService interface {
	Create(ctx context.Context, originalURL string) (string, error)
	Get(ctx context.Context, id string) (string, error)
}
