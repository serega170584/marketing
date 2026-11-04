package handler

import (
	"encoding/json"
	"errors"
	"marketing/internal/infrastructure/database"
	"marketing/internal/infrastructure/telemetry"
	"marketing/internal/usecase"
	"net/http"

	"marketing/internal/transport/http/api"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type UserHandler struct {
	pgConn      *database.Postgres
	userUseCase *usecase.UserUseCase
	logger      *telemetry.Logger
}

func NewUserHandler(
	pgConn *database.Postgres,
	userUseCase *usecase.UserUseCase,
	logger *telemetry.Logger,
) *UserHandler {
	return &UserHandler{
		pgConn:      pgConn,
		userUseCase: userUseCase,
		logger:      logger,
	}
}

func (h *UserHandler) CreateUser(
	w http.ResponseWriter,
	r *http.Request,
) {
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

func (h *UserHandler) GetUsers(
	w http.ResponseWriter,
	r *http.Request,
) {
	users := []api.UserResponse{
		{Id: 1, Name: "John Doe", Email: "john.doe@example.com"},
		{Id: 2, Name: "Jane Doe", Email: "jane.doe@example.com"},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) UpdateUser(
	w http.ResponseWriter,
	r *http.Request,
	id uint64,
) {
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

func (h *UserHandler) GetUserById(
	w http.ResponseWriter,
	r *http.Request,
	id uint64,
) {
	ctx := r.Context()

	user, err := h.userUseCase.GetUserByID(ctx, id)

	if errors.Is(err, usecase.ErrUserNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "User not found"})

		return
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		h.logger.Error(ctx, "getUserByIdHandler: failed to get user by id", map[string]any{
			"error": err,
		})

		_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Internal server error"})

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(api.UserResponse{
		Id:    id,
		Name:  user.Name,
		Email: openapi_types.Email(user.Email),
	})
}
