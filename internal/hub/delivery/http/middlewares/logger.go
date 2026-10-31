package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

var log *slog.Logger

// initLogger добавляет свой логгер в middleware.
func initLogger(loggerSystem *slog.Logger) {
	log = loggerSystem
}

// Logger логирует http запросы.
func Logger(next http.Handler) http.Handler {
	return middleware.RequestLogger(&logger{})(next)
}

type logger struct {
	r *http.Request
}

func (l *logger) NewLogEntry(r *http.Request) middleware.LogEntry {
	l.r = r
	return l
}

func (l *logger) Write(status, bytes int, header http.Header, elapsed time.Duration, extra any) {
	scheme := "http"
	if l.r.TLS != nil {
		scheme = "https"
	}
	path := fmt.Sprintf("%s://%s%s %s", scheme, l.r.Host, l.r.RequestURI, l.r.Proto)

	log.Info(
		"request",
		slog.String("method", l.r.Method),
		slog.String("path", path),
		slog.Int("status", status),
		slog.String("address", l.r.RemoteAddr),
		slog.String("timeExec", elapsed.String()),
		slog.Int("size", bytes),
	)
}

func (l *logger) Panic(v any, stack []byte) {

}
