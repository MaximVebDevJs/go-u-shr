package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	urlService "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	"github.com/MaximVebDevJs/go-u-shr/internal/service/url/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteUrlsService(t *testing.T) {
	t.Parallel()

	const (
		userID  = "user-1"
		baseURL = "http://localhost:8080"
	)

	ctx := context.Background()

	t.Run("дедуплицирует и отправляет batch в репозиторий", func(t *testing.T) {
		t.Parallel()

		repo := mocks.NewUrlRepository(t)
		done := make(chan struct{})

		repo.EXPECT().
			MarkDeleted(mock.Anything, userID, mock.Anything).
			RunAndReturn(func(_ context.Context, gotUserID string, ids []string) error {
				defer close(done)

				assert.Equal(t, userID, gotUserID)
				assert.ElementsMatch(t, []string{"abc12345", "def67890"}, ids)

				return nil
			})

		svc := newURLService(t, repo, baseURL)

		require.NoError(t, svc.DeleteUrls(ctx, userID, []string{"abc12345", "def67890", "abc12345", " "}))

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("MarkDeleted was not called")
		}
	})

	t.Run("большой список режется на чанки", func(t *testing.T) {
		t.Parallel()

		const total = 150

		ids := make([]string, 0, total)
		for i := range total {
			ids = append(ids, fmt.Sprintf("short%04d", i))
		}

		repo := mocks.NewUrlRepository(t)

		var (
			mu       sync.Mutex
			gotSizes []int
			gotIDs   []string
		)

		done := make(chan struct{})

		repo.EXPECT().
			MarkDeleted(mock.Anything, userID, mock.Anything).
			RunAndReturn(func(_ context.Context, _ string, chunk []string) error {
				mu.Lock()
				defer mu.Unlock()

				gotSizes = append(gotSizes, len(chunk))
				gotIDs = append(gotIDs, chunk...)

				if len(gotIDs) == total {
					close(done)
				}

				return nil
			})

		svc := newURLService(t, repo, baseURL)

		require.NoError(t, svc.DeleteUrls(ctx, userID, ids))

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("MarkDeleted did not receive all ids")
		}

		mu.Lock()
		defer mu.Unlock()

		assert.ElementsMatch(t, ids, gotIDs)
		for _, size := range gotSizes {
			assert.LessOrEqual(t, size, 100)
		}
	})

	t.Run("пустой список не доходит до репозитория", func(t *testing.T) {
		t.Parallel()

		repo := mocks.NewUrlRepository(t)
		svc := newURLService(t, repo, baseURL)

		err := svc.DeleteUrls(ctx, userID, []string{})
		require.Error(t, err)
		assert.ErrorIs(t, err, urlService.ErrInvalidUrl)
	})
}
