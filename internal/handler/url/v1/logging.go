package v1

import (
	"net/url"

	"go.uber.org/zap"
)

// urlLogFields возвращает безопасные для логов поля URL без query/fragment (PII redaction).
func urlLogFields(raw string) []zap.Field {
	parsed, err := url.Parse(raw)
	if err != nil {
		return []zap.Field{zap.Int("url_length", len(raw))}
	}

	return []zap.Field{
		zap.String("scheme", parsed.Scheme),
		zap.String("host", parsed.Host),
		zap.Int("url_length", len(raw)),
	}
}
