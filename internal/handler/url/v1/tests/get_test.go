package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	errs "github.com/MaximVebDevJs/go-u-shr/internal/errors"
	apiUrlV1 "github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1"
	"github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUrlHandler(t *testing.T) {
	t.Parallel()

	var (
		ctx         = context.Background()
		id          = "abc123"
		expectedURL = "https://practicum.yandex.ru/"
	)

	tests := []struct {
		name         string
		setupMock    func(svc *mocks.UrlService)
		expectedErr  error
		expectedCode int
		expectedURL  string
	}{
		{
			name: "Успешный редирект",
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().Get(ctx, id).Return(expectedURL, nil)
			},
			expectedErr:  nil,
			expectedCode: http.StatusTemporaryRedirect,
			expectedURL:  expectedURL,
		},
		{
			name: "Урл не найден",
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().Get(ctx, id).Return("", errs.ErrUrlNotFound)
			},
			expectedErr: errs.ErrUrlNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewUrlService(t)

			if tc.setupMock != nil {
				tc.setupMock(svc)
			}

			apiHandler := apiUrlV1.New(svc)

			req := httptest.NewRequest(
				http.MethodGet,
				"/"+id,
				nil,
			)

			req = req.WithContext(ctx)

			req.SetPathValue("id", id)

			rec := httptest.NewRecorder()

			err := apiHandler.GetUrl(rec, req)

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
				return
			}

			require.NoError(t, err)

			assert.Equal(
				t,
				tc.expectedCode,
				rec.Code,
			)

			assert.Equal(
				t,
				tc.expectedURL,
				rec.Header().Get("Location"),
			)
		})
	}
}
