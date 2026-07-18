package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// DecompressMiddleware распаковывает тело запроса, если оно сжато gzip
func DecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "передан невалидный gzip", http.StatusBadRequest)
				return
			}
			defer reader.Close()
			r.Body = io.NopCloser(reader)
			// Удаляем заголовок, чтобы дальше обработчики не знали о сжатии
			r.Header.Del("Content-Encoding")
		}
		next.ServeHTTP(w, r)
	})
}

// CompressMiddleware сжимает ответ, если клиент поддерживает gzip
func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzWriter := &gzipResponseWriter{
			ResponseWriter: w,
			gzWriter:       gz,
			useGzip:        true, // по умолчанию сжимаем
		}
		next.ServeHTTP(gzWriter, r)
	})
}

// gzipResponseWriter перехватывает запись и сжимает только если Content-Type допустим.
type gzipResponseWriter struct {
	http.ResponseWriter
	gzWriter   *gzip.Writer
	useGzip    bool
	headerSent bool
}

// WriteHeader вызывается перед отправкой заголовков, поэтому мы можем принять решение о сжатии.
func (g *gzipResponseWriter) WriteHeader(code int) {
	if !g.headerSent {
		g.headerSent = true
		contentType := g.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") && !strings.Contains(contentType, "text/html") {
			g.useGzip = false
			g.Header().Del("Content-Encoding")
		}
	}
	g.ResponseWriter.WriteHeader(code)
}

// Write переопределяет запись тела, используя соответствующий writer.
func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	if !g.headerSent {
		g.WriteHeader(http.StatusOK)
	}
	if g.useGzip {
		return g.gzWriter.Write(b)
	}
	return g.ResponseWriter.Write(b)
}
