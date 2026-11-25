package middlewares

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery восстанавливает работу программы при панике.
func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				rec := recover()
				if rec == nil {
					return
				}

				// http.ErrAbortHandler — штатный способ прервать запрос, его не логируем
				if rec == http.ErrAbortHandler {
					panic(rec)
				}

				logger.Error("panic recovered",
					slog.Any("panic", rec),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("stack", string(debug.Stack())),
				)

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}()

			next.ServeHTTP(w, r)
		})
	}
}
