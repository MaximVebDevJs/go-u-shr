package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaximVebDevJs/go-u-shr/internal/auth"
	"github.com/MaximVebDevJs/go-u-shr/internal/transport/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureAuthMiddleware(t *testing.T) {
	t.Parallel()

	const userID = "user-1"

	signer := auth.NewSigner("test-secret")

	tests := []struct {
		name        string
		cookieValue string
		wantUserID  string
		wantCookie  bool
	}{
		{
			name:       "cookie отсутствует",
			wantCookie: true,
		},
		{
			name:        "cookie невалидна",
			cookieValue: "invalid-cookie",
			wantCookie:  true,
		},
		{
			name:        "cookie валидна",
			cookieValue: signer.Sign(userID),
			wantUserID:  userID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var gotUserID string
			next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				var ok bool
				gotUserID, ok = auth.UserIDFromContext(r.Context())
				require.True(t, ok)
			})

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			if tc.cookieValue != "" {
				req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: tc.cookieValue})
			}

			rec := httptest.NewRecorder()

			middleware.EnsureAuthMiddleware(signer)(next).ServeHTTP(rec, req)

			if tc.wantUserID != "" {
				assert.Equal(t, tc.wantUserID, gotUserID)
			} else {
				assert.NotEmpty(t, gotUserID)
			}

			cookies := rec.Result().Cookies()
			if !tc.wantCookie {
				assert.Empty(t, cookies)
				return
			}

			require.Len(t, cookies, 1)
			assert.Equal(t, auth.CookieName, cookies[0].Name)

			signedUserID, ok := signer.Verify(cookies[0].Value)
			require.True(t, ok)
			assert.Equal(t, gotUserID, signedUserID)
		})
	}
}

func TestMustAuth(t *testing.T) {
	t.Parallel()

	const userID = "user-1"

	signer := auth.NewSigner("test-secret")

	tests := []struct {
		name         string
		cookieValue  string
		expectedCode int
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
			name:         "cookie валидна",
			cookieValue:  signer.Sign(userID),
			expectedCode: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotUserID, ok := auth.UserIDFromContext(r.Context())
				require.True(t, ok)
				assert.Equal(t, userID, gotUserID)

				w.WriteHeader(http.StatusNoContent)
			})

			req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			if tc.cookieValue != "" {
				req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: tc.cookieValue})
			}

			rec := httptest.NewRecorder()

			middleware.MustAuth(signer)(next).ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
