package http

import (
	"context"
	"marketing/internal/config"
	"marketing/internal/usecase"
	"net/http"

	"marketing/internal/infrastructure/database"
	"marketing/internal/transport/http/api"
	"marketing/internal/transport/http/handler"
	custom_middleware "marketing/internal/transport/http/middleware"

	"github.com/VictoriaMetrics/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type ApiServer struct {
	*handler.UserHandler
	*handler.OrderHandler
}

//nolint:exhaustruct_v5
func newRouter(
	userHandler *handler.UserHandler,
	orderHandler *handler.OrderHandler,
	serviceCfg config.ServiceConfig,
) http.Handler {
	r := chi.NewRouter()

	r.Use(custom_middleware.TraceIDMiddleware)
	r.Use(custom_middleware.LoggerMiddleware())
	r.Use(custom_middleware.AuthMiddleware())
	r.Use(middleware.Recoverer)

	r.Use(func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, serviceCfg.Name)
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	r.Get("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, false)
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
	ctx context.Context,
	pgConn *database.Postgres,
	serviceCfg config.ServiceConfig,
	userUseCase *usecase.UserUseCase,
) http.Handler {
	userHandler := handler.NewUserHandler(
		pgConn,
		userUseCase,
	)

	orderHandler := handler.NewOrderHandler(pgConn)

	return newRouter(
		userHandler,
		orderHandler,
		serviceCfg,
	)
}
