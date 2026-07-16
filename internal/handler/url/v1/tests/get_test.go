package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	errs "github.com/MaximVebDevJs/go-u-shr/internal/errors"
	apiUrlV1 "github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1"
	"github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUrlHandler(t *testing.T) {
	t.Parallel()

	var (
		id          = "abc123"
		expectedURL = "https://practicum.yandex.ru/"
	)

	tests := []struct {
		name         string
		setupMock    func(svc *mocks.UrlService)
		expectedCode int
		expectedURL  string
	}{
		{
			name: "Успешный редирект",
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().Get(mock.Anything, id).Return(expectedURL, nil)
			},
			expectedCode: http.StatusTemporaryRedirect,
			expectedURL:  expectedURL,
		},
		{
			name: "Урл не найден",
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().Get(mock.Anything, id).Return("", errs.ErrUrlNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedURL:  "",
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

			r := chi.NewRouter()
			apiUrlV1.RegisterRoutes(r, apiHandler)

			req := httptest.NewRequest(
				http.MethodGet,
				"/"+id,
				nil,
			)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedURL != "" {
				assert.Equal(t, tc.expectedURL, rec.Header().Get("Location"))
			}
		})
	}
}
