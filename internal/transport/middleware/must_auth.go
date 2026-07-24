package middleware

import (
	"net/http"

	"github.com/MaximVebDevJs/go-u-shr/internal/auth"
)

// MustAuth проверяет наличие валидного пользователя.
// Если пользователь не найден — возвращает 401.
func MustAuth(
	signer *auth.Signer,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			cookieValue, err := auth.GetUserID(r)

			if err != nil {
				http.Error(
					w,
					"Unauthorized",
					http.StatusUnauthorized,
				)
				return
			}

			userID, ok := signer.Verify(cookieValue)

			if !ok {
				http.Error(
					w,
					"Unauthorized",
					http.StatusUnauthorized,
				)
				return
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
