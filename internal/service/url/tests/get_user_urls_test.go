package tests

import (
	"context"
	"testing"

	"github.com/MaximVebDevJs/go-u-shr/internal/model"
	urlService "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	"github.com/MaximVebDevJs/go-u-shr/internal/service/url/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserURLsService(t *testing.T) {
	t.Parallel()

	const (
		userID  = "user-1"
		baseURL = "http://localhost:8080"
	)

	ctx := context.Background()
	repo := mocks.NewUrlRepository(t)
	repo.EXPECT().
		GetUserURLs(ctx, userID).
		Return([]model.UserURL{
			{ShortID: "abc12345", OriginalURL: "https://example.com/a"},
		}, nil)

	svc := urlService.New(repo, baseURL)

	got, err := svc.GetUserURLs(ctx, userID)
	require.NoError(t, err)

	require.Len(t, got, 1)
	assert.Equal(t, baseURL+"/abc12345", got[0].ShortURL)
	assert.Equal(t, "https://example.com/a", got[0].OriginalURL)
}
