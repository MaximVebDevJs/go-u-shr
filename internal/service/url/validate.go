package url

import (
	"net/url"
	"strings"
)

// validateURL проверяет, что строка — допустимый http/https URL.
// Отклоняем javascript:, data:, protocol-relative (//) и пустые значения,
// чтобы shortener не использовался как open redirect.
func validateURL(raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ErrInvalidUrl
	}

	u, err := url.Parse(raw)
	if err != nil {
		return ErrInvalidUrl
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return ErrInvalidUrl
	}

	if u.Host == "" {
		return ErrInvalidUrl
	}

	return nil
}
