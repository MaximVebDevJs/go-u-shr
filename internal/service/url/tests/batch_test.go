package tests

import (
	"context"
	"strings"
	"testing"

	urlRepo "github.com/MaximVebDevJs/go-u-shr/internal/repository/url"
	urlService "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	"github.com/MaximVebDevJs/go-u-shr/internal/service/url/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBatchCreateUrlService(t *testing.T) {
	t.Parallel()

	var (
		ctx     = context.Background()
		userID  = "user-1"
		baseURL = "http://localhost:8080"
	)

	tests := []struct {
		name               string
		items              []urlService.BatchItem
		setupMock          func(repo *mocks.UrlRepository)
		expectedErr        error
		checkDuplicateURLs []string
	}{
		{
			name: "успешное пакетное создание коротких ссылок",
			items: []urlService.BatchItem{
				{CorrelationID: "1", OriginalURL: "https://example.com/a"},
				{CorrelationID: "2", OriginalURL: "https://example.com/b"},
			},
			setupMock: func(repo *mocks.UrlRepository) {
				repo.EXPECT().
					CreateBatch(ctx, userID, mock.Anything).
					RunAndReturn(func(_ context.Context, gotUserID string, records []urlRepo.BatchRecord) error {
						assert.Equal(t, userID, gotUserID)
						require.Len(t, records, 2)
						assert.Equal(t, "https://example.com/a", records[0].OriginalURL)
						assert.Equal(t, "https://example.com/b", records[1].OriginalURL)
						assert.NotEmpty(t, records[0].ID)
						assert.NotEmpty(t, records[1].ID)
						assert.NotEqual(t, records[0].ID, records[1].ID)

						return nil
					})
			},
		},
		{
			name:  "пустой batch не обращается к репозиторию",
			items: []urlService.BatchItem{},
			setupMock: func(repo *mocks.UrlRepository) {
			},
		},
		{
			name: "невалидный url не доходит до репозитория",
			items: []urlService.BatchItem{
				{CorrelationID: "1", OriginalURL: "javascript:alert(1)"},
			},
			setupMock:   func(repo *mocks.UrlRepository) {},
			expectedErr: urlService.ErrInvalidUrl,
		},
		{
			name: "дубликат original_url внутри batch",
			items: []urlService.BatchItem{
				{CorrelationID: "1", OriginalURL: "https://example.com/a"},
				{CorrelationID: "2", OriginalURL: "https://example.com/a"},
			},
			setupMock: func(repo *mocks.UrlRepository) {},
			checkDuplicateURLs: []string{
				"https://example.com/a",
			},
		},
		{
			name: "превышен лимит batch",
			items: func() []urlService.BatchItem {
				items := make([]urlService.BatchItem, 101)
				for i := range items {
					items[i] = urlService.BatchItem{
						CorrelationID: "id",
						OriginalURL:   "https://example.com/page",
					}
				}

				return items
			}(),
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

			svc := newURLService(t, repo, baseURL)

			results, err := svc.BatchCreate(ctx, userID, tc.items)

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, results)

				return
			}

			if tc.checkDuplicateURLs != nil {
				require.Error(t, err)

				var duplicateErr *urlService.DuplicateURLsError
				require.ErrorAs(t, err, &duplicateErr)
				assert.ElementsMatch(t, tc.checkDuplicateURLs, duplicateErr.URLs)
				assert.Nil(t, results)

				return
			}

			require.NoError(t, err)
			require.Len(t, results, len(tc.items))

			for i, result := range results {
				assert.Equal(t, tc.items[i].CorrelationID, result.CorrelationID)
				assert.True(t, strings.HasPrefix(result.ShortURL, baseURL+"/"))

				id := strings.TrimPrefix(result.ShortURL, baseURL+"/")
				assert.Len(t, id, 8)
			}
		})
	}
}
