package auth

import "net/http"

// CookieName содержит имя cookie с подписанным ID пользователя.
const CookieName = "user_id"

// GetUserID возвращает сырое подписанное значение user_id cookie из запроса.
func GetUserID(r *http.Request) (string, error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

// SetUserID устанавливает cookie с подписанным ID пользователя.
func SetUserID(w http.ResponseWriter, signedUserID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    signedUserID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
