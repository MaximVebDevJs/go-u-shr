package middleware

import (
	"net/http"

	"github.com/MaximVebDevJs/go-u-shr/internal/auth"
	"github.com/google/uuid"
)

// EnsureAuthMiddleware гарантирует наличие пользователя.
// Если cookie отсутствует или невалидна — создаёт нового пользователя.
func EnsureAuthMiddleware(
	signer *auth.Signer,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			var userID string

			cookieValue, err := auth.GetUserID(r)

			if err == nil {
				id, ok := signer.Verify(cookieValue)

				if ok {
					userID = id
				}
			}

			// Cookie нет или подпись невалидна
			if userID == "" {

				userID = uuid.NewString()

				signedValue := signer.Sign(userID)

				auth.SetUserID(
					w,
					signedValue,
				)
			}

			ctx := auth.WithUserID(
				r.Context(),
				userID,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}
