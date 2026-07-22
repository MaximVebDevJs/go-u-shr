package v1

import (
	"context"

	serviceurl "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
)

type UrlService interface {
	Create(ctx context.Context, originalURL string) (string, error)
	BatchCreate(ctx context.Context, items []serviceurl.BatchItem) ([]serviceurl.BatchResult, error)
	Get(ctx context.Context, id string) (string, error)
	Ping(ctx context.Context) error
}
