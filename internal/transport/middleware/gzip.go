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
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "передан невалидный gzip", http.StatusBadRequest)
			return
		}
		defer reader.Close()

		r.Body = io.NopCloser(reader)
		r.Header.Del("Content-Encoding")

		next.ServeHTTP(w, r)
	})
}

// CompressMiddleware сжимает ответ
func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzw := &gzipResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		defer func() {
			if gzw.gzipWriter != nil {
				_ = gzw.gzipWriter.Close()
			}
		}()

		next.ServeHTTP(gzw, r)
	})
}

// gzipResponseWriter решает,
// нужно ли сжимать ответ.
type gzipResponseWriter struct {
	http.ResponseWriter

	gzipWriter *gzip.Writer

	statusCode int

	headerReceived bool
	headerSent     bool
	compressed     bool
}

// WriteHeader только запоминает статус.
// Реальные заголовки отправляются при первом Write().
func (g *gzipResponseWriter) WriteHeader(code int) {
	if g.headerReceived {
		return
	}

	g.statusCode = code
	g.headerReceived = true
}

// Write отправляет заголовки и тело.
func (g *gzipResponseWriter) Write(b []byte) (int, error) {

	// Если обработчик не указал Content-Type,
	// определяем его автоматически.
	if g.Header().Get("Content-Type") == "" {
		g.Header().Set("Content-Type", http.DetectContentType(b))
	}

	if !g.compressed && shouldCompress(g.Header().Get("Content-Type")) {
		g.compressed = true

		g.Header().Set("Content-Encoding", "gzip")
		g.Header().Del("Content-Length")

		g.gzipWriter = gzip.NewWriter(g.ResponseWriter)
	}

	if !g.headerSent {
		g.headerSent = true
		g.ResponseWriter.WriteHeader(g.statusCode)
	}

	if g.compressed {
		return g.gzipWriter.Write(b)
	}

	return g.ResponseWriter.Write(b)
}

// Flush поддерживает потоковую отправку.
func (g *gzipResponseWriter) Flush() {
	if g.gzipWriter != nil {
		_ = g.gzipWriter.Flush()
	}

	if flusher, ok := g.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// shouldCompress определяет,
// нужно ли сжимать данный тип контента.
func shouldCompress(contentType string) bool {
	return strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "text/html")
}
