package tests

import (
	"context"
	"testing"
	"time"

	"github.com/MaximVebDevJs/go-u-shr/internal/logger"
	urlService "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
)

type urlServiceForTests interface {
	Create(ctx context.Context, userID string, originalURL string) (string, error)
	BatchCreate(ctx context.Context, userID string, items []urlService.BatchItem) ([]urlService.BatchResult, error)
	Get(ctx context.Context, id string) (string, error)
	GetUserURLs(ctx context.Context, userID string) ([]urlService.UserURL, error)
	DeleteUrls(ctx context.Context, userID string, ids []string) error
	Close(ctx context.Context) error
}

func newURLService(t *testing.T, repo urlService.UrlRepository, baseURL string) urlServiceForTests {
	t.Helper()

	svc := urlService.New(repo, baseURL, logger.Nop())
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := svc.Close(ctx); err != nil {
			t.Fatalf("close url service: %v", err)
		}
	})

	return svc
}
