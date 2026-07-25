package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MaximVebDevJs/go-u-shr/internal/auth"
	apiUrlV1 "github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1"
	"github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1/mocks"
	"github.com/MaximVebDevJs/go-u-shr/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteUrlsHandler(t *testing.T) {
	t.Parallel()

	const userID = "user-1"

	signer := auth.NewSigner("test-secret")

	tests := []struct {
		name         string
		cookieValue  string
		body         string
		setupMock    func(svc *mocks.UrlService)
		expectedCode int
	}{
		{
			name:         "cookie отсутствует",
			body:         `["abc12345"]`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "битый JSON",
			cookieValue:  signer.Sign(userID),
			body:         `[`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "пустой список",
			cookieValue:  signer.Sign(userID),
			body:         `[]`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "запрос принят",
			cookieValue: signer.Sign(userID),
			body:        `["abc12345","def67890"]`,
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().
					DeleteUrls(mock.Anything, userID, []string{"abc12345", "def67890"}).
					Return(nil)
			},
			expectedCode: http.StatusAccepted,
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

			req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(tc.body))
			if tc.cookieValue != "" {
				req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: tc.cookieValue})
			}

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
