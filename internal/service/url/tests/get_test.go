package tests

import (
	"context"
	"testing"

	"github.com/MaximVebDevJs/go-u-shr/internal/model"
	errs "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	"github.com/MaximVebDevJs/go-u-shr/internal/service/url/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUrlService(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()
		id  = "EwHXdJfB"

		expectedURL = "https://practicum.yandex.ru/"
		baseURL     = "http://localhost:8080"
	)

	tests := []struct {
		name        string
		id          string
		setupMock   func(repo *mocks.UrlRepository)
		expectedURL string
		expectedErr error
	}{
		{
			name: "успешное получение URL",
			id:   id,
			setupMock: func(repo *mocks.UrlRepository) {
				repo.EXPECT().
					Get(ctx, id).
					Return(expectedURL, nil)
			},
			expectedURL: expectedURL,
			expectedErr: nil,
		},
		{
			name: "url не найден в репозитории",
			id:   id,
			setupMock: func(repo *mocks.UrlRepository) {
				repo.EXPECT().
					Get(ctx, id).
					Return("", model.ErrURLNotFound)
			},
			expectedURL: "",
			expectedErr: errs.ErrUrlNotFound,
		},
		{
			name: "url удалён в репозитории",
			id:   id,
			setupMock: func(repo *mocks.UrlRepository) {
				repo.EXPECT().
					Get(ctx, id).
					Return("", model.ErrURLDeleted)
			},
			expectedURL: "",
			expectedErr: errs.ErrURLDeleted,
		},
		{
			name: "ошибка получения URL из репозитория (legacy mock)",
			id:   id,
			setupMock: func(repo *mocks.UrlRepository) {
				repo.EXPECT().
					Get(ctx, id).
					Return("", errs.ErrUrlNotFound)
			},
			expectedURL: "",
			expectedErr: errs.ErrUrlNotFound,
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

			result, err := svc.Get(ctx, tc.id)

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)

			assert.Equal(
				t,
				tc.expectedURL,
				result,
			)
		})
	}
}
