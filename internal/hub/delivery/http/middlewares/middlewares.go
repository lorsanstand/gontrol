package middlewares

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
)

// ConnectMiddlewares подключает middleware к основному роуту.
func ConnectMiddlewares(r *chi.Mux, logger *slog.Logger) {
	initLogger(logger)

	r.Use(Logger)
	r.Use(Recovery)
}
