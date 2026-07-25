package v1

import (
	"context"

	serviceurl "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
)

type UrlService interface {
	Create(ctx context.Context, userID string, originalURL string) (string, error)
	BatchCreate(ctx context.Context, userID string, items []serviceurl.BatchItem) ([]serviceurl.BatchResult, error)
	Get(ctx context.Context, id string) (string, error)
	GetUserURLs(ctx context.Context, userID string) ([]serviceurl.UserURL, error)
	DeleteUrls(ctx context.Context, userID string, ids []string) error
	Ping(ctx context.Context) error
}
