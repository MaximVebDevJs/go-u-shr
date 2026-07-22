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

	urlRepo "github.com/MaximVebDevJs/go-u-shr/internal/repository/url"
)

func TestCreateUrlService(t *testing.T) {
	t.Parallel()

	var (
		ctx         = context.Background()
		originalURL = "https://practicum.yandex.ru/"
		baseURL     = "http://localhost:8080"
	)

	tests := []struct {
		name           string
		inputURL       string
		setupMock      func(repo *mocks.UrlRepository)
		expectedErr    error
		expectedResult string
	}{
		{
			name:     "успешное создание короткой ссылки",
			inputURL: originalURL,
			setupMock: func(repo *mocks.UrlRepository) {
				repo.EXPECT().
					Create(
						ctx,
						originalURL,
						mock.Anything,
					).
					Return("abc12345", nil)
			},
		},
		{
			name:     "url уже существует",
			inputURL: originalURL,
			setupMock: func(repo *mocks.UrlRepository) {
				repo.EXPECT().
					Create(
						ctx,
						originalURL,
						mock.Anything,
					).
					Return("existing1", urlRepo.ErrOriginalURLExists)
			},
			expectedErr:    urlService.ErrURLAlreadyExists,
			expectedResult: baseURL + "/existing1",
		},
		{
			name:        "невалидный url не доходит до репозитория",
			inputURL:    "javascript:alert(1)",
			setupMock:   func(repo *mocks.UrlRepository) {},
			expectedErr: urlService.ErrInvalidUrl,
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
				tc.inputURL,
			)

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Equal(t, tc.expectedResult, result)

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
