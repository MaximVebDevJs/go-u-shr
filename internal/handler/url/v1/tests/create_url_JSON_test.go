package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	serviceurl "github.com/MaximVebDevJs/go-u-shr/internal/service/url"

	apiUrlV1 "github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1"
	"github.com/MaximVebDevJs/go-u-shr/internal/handler/url/v1/mocks"
	"github.com/MaximVebDevJs/go-u-shr/internal/logger"
	transporthttp "github.com/MaximVebDevJs/go-u-shr/internal/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUrlJSONHandler(t *testing.T) {
	t.Parallel()

	var (
		ctx         = context.Background()
		originalURL = "https://practicum.yandex.ru/"
		shortURL    = "http://localhost:8080/abc123"
	)

	tests := []struct {
		name         string
		body         interface{}
		setupMock    func(svc *mocks.UrlService)
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name: "успешное создание короткой ссылки",
			body: map[string]string{"url": originalURL},
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().
					Create(ctx, originalURL).
					Return(shortURL, nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: map[string]interface{}{"result": shortURL},
		},
		{
			name:         "передан пустой URL",
			body:         map[string]string{"url": ""},
			setupMock:    func(svc *mocks.UrlService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": serviceurl.ErrInvalidUrl.Error(),
				"code":    float64(http.StatusBadRequest),
			},
		},
		{
			name:         "передан невалидный JSON",
			body:         `{"url": "missing closing brace"`,
			setupMock:    func(svc *mocks.UrlService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": serviceurl.ErrInvalidJSON.Error(),
				"code":    float64(http.StatusBadRequest),
			},
		},
		{
			name: "ошибка сервиса",
			body: map[string]string{"url": originalURL},
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().
					Create(ctx, originalURL).
					Return("", errors.New("internal error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"message": "внутренняя ошибка",
				"code":    float64(http.StatusInternalServerError),
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

			var reqBody []byte
			switch v := tc.body.(type) {
			case map[string]string:
				reqBody, _ = json.Marshal(v)
			case string:
				reqBody = []byte(v)
			default:
				reqBody = []byte{}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()

			handler := transporthttp.Wrap(apiHandler.CreateUrlJSON, transporthttp.ErrorHandler, log)
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			if tc.expectedBody != nil {
				var resp map[string]interface{}
				errDecode := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, errDecode)
				for key, val := range tc.expectedBody {
					assert.Equal(t, val, resp[key])
				}
			}
		})
	}
}
