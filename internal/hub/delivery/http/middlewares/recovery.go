package middlewares

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// Recovery восстанавливает работу программы при панике.
func Recovery(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}
