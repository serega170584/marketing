package handler

import (
	"encoding/json"
	"marketing/internal/infrastructure/database"
	"net/http"

	"marketing/internal/transport/http/api"
)

type OrderHandler struct {
	pgConn *database.Postgres
}

func NewOrderHandler(pgConn *database.Postgres) *OrderHandler {
	return &OrderHandler{
		pgConn: pgConn,
	}
}

// CreateOrder обрабатывает POST /orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req api.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid request payload"})

		return
	}

	// Здесь в будущем будет вызов h.pgConn.Pool.Exec(...) для сохранения в БД

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(api.CreateOrderResponse{
		Status:  "success",
		Message: "Order created",
	})
}

// GetOrders обрабатывает GET /orders
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	// В будущем здесь будет получение данных из БД через h.pgConn
	orders := []api.OrderResponse{
		{Id: 42, UserId: 1, Status: "pending", TotalAmount: 99.99},
		{Id: 43, UserId: 2, Status: "processing", TotalAmount: 145.50},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(orders)
}

// UpdateOrder обрабатывает PUT /orders/{id}
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request, id int) {
	var req api.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid request payload"})

		return
	}

	// Заглушка для проверки 404 (Not Found), аналогичная вашей в UserHandler
	if id != 42 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Order not found"})

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(api.OrderResponse{
		Id:          id,
		UserId:      1, // В реальном кейсе здесь будут актуальные данные из БД
		Status:      req.Status,
		TotalAmount: req.TotalAmount,
	})
}
