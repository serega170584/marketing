package http

import (
	"net/http"

	"marketing/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(userHandler *handler.UserHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/users", userHandler.Create)
		r.Get("/users", userHandler.Get)
	})

	return r
}
