package tests

import (
	"context"
	"strings"
	"testing"

	urlService "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	"github.com/MaximVebDevJs/go-u-shr/internal/service/url/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUrlService(t *testing.T) {
	t.Parallel()

	var (
		ctx         = context.Background()
		originalURL = "https://practicum.yandex.ru/"
		baseURL     = "http://localhost:8080"
	)

	tests := []struct {
		name        string
		setupMock   func(repo *mocks.UrlRepository)
		expectedErr error
	}{
		{
			name: "успешное создание короткой ссылки",
			setupMock: func(repo *mocks.UrlRepository) {
				repo.EXPECT().
					Create(
						ctx,
						originalURL,
						mock.Anything,
					).
					Return(nil)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewUrlRepository(t)

			if tc.setupMock != nil {
				tc.setupMock(repo)
			}

			svc := urlService.New(repo, baseURL)

			result, err := svc.Create(
				ctx,
				originalURL,
			)

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Empty(t, result)
				return
			}

			require.NoError(t, err)

			assert.True(
				t,
				strings.HasPrefix(result, baseURL+"/"),
			)

			id := strings.TrimPrefix(
				result,
				baseURL+"/",
			)

			assert.Len(
				t,
				id,
				8,
			)
		})
	}
}
