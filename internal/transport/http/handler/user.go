package handler

import "net/http"

type UserHandler struct {
	// Здесь в будущем будет зависимость от бизнес-логики: userService *service.UserService
}

// NewUserHandler — конструктор хендлера
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// Create обрабатывает POST-запрос на создание пользователя
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Имитация успешного ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(`{"status": "success", "message": "User created"}`))
}

// Get обрабатывает GET-запрос на получение данных
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"id": 1, "name": "John Doe"}`))
}
