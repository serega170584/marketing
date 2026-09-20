package handler

import (
	"encoding/json"
	"net/http"

	"marketing/internal/transport/http/api"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req api.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid request payload"})

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(api.CreateUserResponse{
		Status:  "success",
		Message: "User created",
	})
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := []api.UserResponse{
		{Id: 1, Name: "John Doe", Email: "john.doe@example.com"},
		{Id: 2, Name: "Jane Doe", Email: "jane.doe@example.com"},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, id int) {
	var req api.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid request payload"})

		return
	}

	if id != 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound) // Fixed: removed the undefined slice header syntax
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "User not found"})

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(api.UserResponse{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	})
}
