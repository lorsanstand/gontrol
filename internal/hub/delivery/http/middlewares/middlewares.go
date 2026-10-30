package middlewares

import (
	"github.com/go-chi/chi/v5"
)

func ConnectMiddlewares(r *chi.Mux) {
	r.Use(Logger)
	r.Use(Recovery)
}