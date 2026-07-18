package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// DecompressMiddleware распаковывает тело запроса,
// если клиент прислал Content-Encoding: gzip.
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

// CompressMiddleware подменяет ResponseWriter,
// если клиент поддерживает gzip.
func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzw := &gzipResponseWriter{
			ResponseWriter: w,
		}
		defer func() {
			if gzw.gzipWriter != nil {
				_ = gzw.gzipWriter.Close()
			}
		}()

		next.ServeHTTP(gzw, r)
	})
}

// gzipResponseWriter сжимает только JSON и HTML.
type gzipResponseWriter struct {
	http.ResponseWriter

	gzipWriter *gzip.Writer
	compressed bool
}

// WriteHeader вызывается обработчиком.
// Если Content-Type уже известен и подходит,
// включаем gzip до отправки заголовков.
func (g *gzipResponseWriter) WriteHeader(code int) {
	if shouldCompress(g.Header().Get("Content-Type")) {
		g.enableCompression()
	}

	g.ResponseWriter.WriteHeader(code)
}

// Write записывает тело ответа.
// Если обработчик не указал Content-Type,
// определяем его автоматически.
func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	if !g.compressed {
		if g.Header().Get("Content-Type") == "" {
			g.Header().Set("Content-Type", http.DetectContentType(b))
		}

		if shouldCompress(g.Header().Get("Content-Type")) {
			g.enableCompression()
		}
	}

	if g.compressed {
		return g.gzipWriter.Write(b)
	}

	return g.ResponseWriter.Write(b)
}

// Flush поддерживает потоковую отправку ответа.
func (g *gzipResponseWriter) Flush() {
	if g.compressed {
		_ = g.gzipWriter.Flush()
	}

	if flusher, ok := g.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// enableCompression включает gzip один раз.
func (g *gzipResponseWriter) enableCompression() {
	if g.compressed {
		return
	}

	g.compressed = true
	g.Header().Set("Content-Encoding", "gzip")
	g.Header().Del("Content-Length")
	g.gzipWriter = gzip.NewWriter(g.ResponseWriter)
}

// shouldCompress определяет,
// нужно ли сжимать данный Content-Type.
func shouldCompress(contentType string) bool {
	return strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "text/html")
}
