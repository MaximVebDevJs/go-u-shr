package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apiUrlV1 "github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1"
	"github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1/mocks"
	"github.com/MaximVebDevJs/go-u-shr/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPingHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupMock    func(svc *mocks.UrlService)
		expectedCode int
		expectJSON   bool
	}{
		{
			name: "база данных доступна",
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().Ping(mock.Anything).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "база данных недоступна",
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().Ping(mock.Anything).Return(errors.New("connection refused"))
			},
			expectedCode: http.StatusInternalServerError,
			expectJSON:   true,
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
			apiUrlV1.RegisterRoutes(r, apiHandler, log)

			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectJSON {
				var resp map[string]interface{}
				err := json.NewDecoder(rec.Body).Decode(&resp)

				assert.NoError(t, err)
				assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
				assert.Equal(t, "внутренняя ошибка", resp["message"])
				assert.Equal(t, float64(http.StatusInternalServerError), resp["code"])
			}
		})
	}
}
