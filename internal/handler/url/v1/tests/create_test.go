package tests

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	apiUrlV1 "github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1"
	"github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1/mocks"
	"github.com/MaximVebDevJs/go-u-shr/internal/logger"
	serviceurl "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUrlHandler(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()

		originalURL = "https://practicum.yandex.ru/"
		shortURL    = "http://localhost:8080/EwHXdJfB"
	)

	tests := []struct {
		name         string
		body         string
		setupMock    func(svc *mocks.UrlService)
		expectedErr  error
		expectedCode int
		expectedBody string
		expectedType string
	}{
		{
			name: "успешное создание короткой ссылки",
			body: originalURL,
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().
					Create(ctx, originalURL).
					Return(shortURL, nil)
			},
			expectedErr:  nil,
			expectedCode: http.StatusCreated,
			expectedBody: shortURL,
			expectedType: "text/plain",
		},
		{
			name: "url уже существует",
			body: originalURL,
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().
					Create(ctx, originalURL).
					Return(shortURL, serviceurl.ErrURLAlreadyExists)
			},
			expectedErr:  nil,
			expectedCode: http.StatusConflict,
			expectedBody: shortURL,
			expectedType: "text/plain",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewUrlService(t)

			if tc.setupMock != nil {
				tc.setupMock(svc)
			}

			apiHandler := apiUrlV1.New(svc, logger.Nop())

			req := httptest.NewRequest(
				http.MethodPost,
				"/",
				bytes.NewBufferString(tc.body),
			)

			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()

			err := apiHandler.CreateUrl(rec, req)

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
				tc.expectedType,
				rec.Header().Get("Content-Type"),
			)

			assert.Equal(
				t,
				tc.expectedBody,
				rec.Body.String(),
			)
		})
	}
}
