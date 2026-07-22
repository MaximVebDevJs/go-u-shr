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

func TestBatchUrlsHandler(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name         string
		body         interface{}
		setupMock    func(svc *mocks.UrlService)
		expectedCode int
		expectedBody interface{}
	}{
		{
			name: "успешное пакетное создание коротких ссылок",
			body: []map[string]string{
				{"correlation_id": "1", "original_url": "https://example.com/a"},
				{"correlation_id": "2", "original_url": "https://example.com/b"},
			},
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().
					BatchCreate(ctx, []serviceurl.BatchItem{
						{CorrelationID: "1", OriginalURL: "https://example.com/a"},
						{CorrelationID: "2", OriginalURL: "https://example.com/b"},
					}).
					Return([]serviceurl.BatchResult{
						{CorrelationID: "1", ShortURL: "http://localhost:8080/abc12345"},
						{CorrelationID: "2", ShortURL: "http://localhost:8080/def67890"},
					}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: []map[string]string{
				{"correlation_id": "1", "short_url": "http://localhost:8080/abc12345"},
				{"correlation_id": "2", "short_url": "http://localhost:8080/def67890"},
			},
		},
		{
			name: "пустой массив",
			body: []map[string]string{},
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().
					BatchCreate(ctx, []serviceurl.BatchItem{}).
					Return([]serviceurl.BatchResult{}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: []map[string]string{},
		},
		{
			name: "передан пустой correlation_id",
			body: []map[string]string{
				{"correlation_id": "", "original_url": "https://example.com/a"},
			},
			setupMock:    func(svc *mocks.UrlService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": serviceurl.ErrInvalidUrl.Error(),
				"code":    float64(http.StatusBadRequest),
			},
		},
		{
			name: "передан пустой original_url",
			body: []map[string]string{
				{"correlation_id": "1", "original_url": ""},
			},
			setupMock:    func(svc *mocks.UrlService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": serviceurl.ErrInvalidUrl.Error(),
				"code":    float64(http.StatusBadRequest),
			},
		},
		{
			name:         "передан невалидный JSON",
			body:         `[{"correlation_id": "1"`,
			setupMock:    func(svc *mocks.UrlService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": serviceurl.ErrInvalidJSON.Error(),
				"code":    float64(http.StatusBadRequest),
			},
		},
		{
			name: "ошибка сервиса",
			body: []map[string]string{
				{"correlation_id": "1", "original_url": "https://example.com/a"},
			},
			setupMock: func(svc *mocks.UrlService) {
				svc.EXPECT().
					BatchCreate(ctx, []serviceurl.BatchItem{
						{CorrelationID: "1", OriginalURL: "https://example.com/a"},
					}).
					Return(nil, errors.New("internal error"))
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
			case []map[string]string:
				reqBody, _ = json.Marshal(v)
			case string:
				reqBody = []byte(v)
			default:
				reqBody = []byte{}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()

			handler := transporthttp.Wrap(apiHandler.BatchUrls, transporthttp.ErrorHandler, log)
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			switch expected := tc.expectedBody.(type) {
			case map[string]interface{}:
				var resp map[string]interface{}
				errDecode := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, errDecode)
				for key, val := range expected {
					assert.Equal(t, val, resp[key])
				}
			case []map[string]string:
				var resp []map[string]string
				errDecode := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, errDecode)
				assert.Equal(t, expected, resp)
			}
		})
	}
}
