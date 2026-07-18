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

// CompressMiddleware запаковывает запрос
func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzw := &gzipResponseWriter{
			ResponseWriter: w,
			gzipWriter:     gzip.NewWriter(w),
		}
		defer gzw.gzipWriter.Close()

		next.ServeHTTP(gzw, r)
	})
}

// gzipResponseWriter перехватывает запись ответа
// и решает, нужно ли использовать gzip.
type gzipResponseWriter struct {
	http.ResponseWriter

	gzipWriter  *gzip.Writer
	compressed  bool
	headerWrote bool
}

// Write вызывается всеми обработчиками при записи ответа.
// Здесь определяется Content-Type и принимается решение,
// сжимать ответ или нет.
func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	if !g.headerWrote {
		g.headerWrote = true

		// Если обработчик сам не указал Content-Type,
		// определяем его автоматически.
		if g.Header().Get("Content-Type") == "" {
			g.Header().Set("Content-Type", http.DetectContentType(b))
		}

		if shouldCompress(g.Header().Get("Content-Type")) {
			g.compressed = true
			g.Header().Set("Content-Encoding", "gzip")
			g.Header().Del("Content-Length")
		}

		g.ResponseWriter.WriteHeader(http.StatusOK)
	}

	if g.compressed {
		return g.gzipWriter.Write(b)
	}

	return g.ResponseWriter.Write(b)
}

// WriteHeader нужен, чтобы корректно обработать случаи,
// когда обработчик сам вызывает w.WriteHeader(...)
// до записи тела ответа.
func (g *gzipResponseWriter) WriteHeader(statusCode int) {
	if g.headerWrote {
		return
	}

	g.headerWrote = true
	g.ResponseWriter.WriteHeader(statusCode)
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

// shouldCompress определяет,
// нужно ли сжимать ответ данного Content-Type.
func shouldCompress(contentType string) bool {
	return strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "text/html")
}
