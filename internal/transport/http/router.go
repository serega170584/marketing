package http

import (
	"net/http"

	"marketing/internal/transport/http/api"
	"marketing/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//nolint:exhaustruct_v5
func NewRouter(userHandler *handler.UserHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	return api.HandlerWithOptions(userHandler, api.ChiServerOptions{
		BaseRouter: r,
		BaseURL:    "/api/v1",
	})
}
