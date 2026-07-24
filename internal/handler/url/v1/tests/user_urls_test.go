package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaximVebDevJs/go-u-shr/internal/auth"
	apiUrlV1 "github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1"
	"github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1/mocks"
	"github.com/MaximVebDevJs/go-u-shr/internal/logger"
	serviceurl "github.com/MaximVebDevJs/go-u-shr/internal/service/url"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetUserURLsHandler(t *testing.T) {
	t.Parallel()

	const userID = "user-1"

	signer := auth.NewSigner("test-secret")

	tests := []struct {
		name         string
		cookieValue  string
		setupMock    func(svc *mocks.UrlService)
		expectedCode int
		expectedBody []map[string]string
	}{
		{
			name:         "cookie отсутствует",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "cookie невалидна",
			cookieValue:  "invalid-cookie",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:        "у пользователя нет ссылок",
			cookieValue: signer.Sign(userID),
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().
					GetUserURLs(mock.Anything, userID).
					Return([]serviceurl.UserURL{}, nil)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:        "ссылки пользователя возвращаются",
			cookieValue: signer.Sign(userID),
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().
					GetUserURLs(mock.Anything, userID).
					Return([]serviceurl.UserURL{
						{
							ShortURL:    "http://localhost:8080/abc12345",
							OriginalURL: "https://example.com/a",
						},
					}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: []map[string]string{
				{
					"short_url":    "http://localhost:8080/abc12345",
					"original_url": "https://example.com/a",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewUrlService(t)
			if tc.setupMock != nil {
				tc.setupMock(svc)
			}

			log := logger.Nop()
			apiHandler := apiUrlV1.New(svc, log)

			r := chi.NewRouter()
			apiUrlV1.RegisterRoutes(r, apiHandler, log, signer)

			req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			if tc.cookieValue != "" {
				req.AddCookie(&http.Cookie{
					Name:  auth.CookieName,
					Value: tc.cookieValue,
				})
			}

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedBody == nil {
				return
			}

			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var got []map[string]string
			err := json.NewDecoder(rec.Body).Decode(&got)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedBody, got)
		})
	}
}
