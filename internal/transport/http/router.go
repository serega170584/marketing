package http

import (
	"marketing/internal/infrastructure/database"
	"net/http"

	"marketing/internal/transport/http/api"
	"marketing/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ApiServer struct {
	*handler.UserHandler
	*handler.OrderHandler
}

//nolint:exhaustruct_v5
func newRouter(
	userHandler *handler.UserHandler,
	orderHandler *handler.OrderHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	apiServer := &ApiServer{
		UserHandler:  userHandler,
		OrderHandler: orderHandler,
	}

	return api.HandlerWithOptions(apiServer, api.ChiServerOptions{
		BaseRouter: r,
		BaseURL:    "/api/v1",
	})
}

func GetRouter(
	pgConn *database.Postgres,
) http.Handler {
	userHandler := handler.NewUserHandler(pgConn)
	orderHandler := handler.NewOrderHandler(pgConn)

	return newRouter(userHandler, orderHandler)
}
