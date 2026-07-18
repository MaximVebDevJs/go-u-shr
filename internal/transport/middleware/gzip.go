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

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")

		gzWriter := &gzipResponseWriter{
			ResponseWriter: w,
			Writer:         gz,
		}
		next.ServeHTTP(gzWriter, r)
	})
}

// gzipResponseWriter перехватывает запись и сжимает только если Content-Type допустим.
type gzipResponseWriter struct {
	http.ResponseWriter
	Writer  io.Writer
	written bool
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	if !g.written {
		contentType := g.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") &&
			!strings.Contains(contentType, "text/html") {
			g.Header().Del("Content-Encoding")
			return g.ResponseWriter.Write(b)
		}
		g.written = true
	}
	return g.Writer.Write(b)
}
