package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"marketing/internal/infrastructure/database"
	"marketing/internal/transport/http/api"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type PromotionHandler struct {
	pgConn *database.Postgres
}

func NewPromotionHandler(pgConn *database.Postgres) *PromotionHandler {
	return &PromotionHandler{
		pgConn: pgConn,
	}
}

func (h *PromotionHandler) CreatePromotion(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req api.CreatePromotionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{
			Error: "Invalid request payload",
		})

		return
	}

	response := api.PromotionResponse{
		Id:            uuid.New(),
		Name:          req.Name,
		DiscountType:  req.DiscountType,
		DiscountValue: req.DiscountValue,
		MaxDiscount:   req.MaxDiscount,
		MaxUses:       req.MaxUses,
		UsedCount:     0,
		StartsAt:      req.StartsAt,
		ExpiresAt:     req.ExpiresAt,
		Status:        req.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *PromotionHandler) GetPromotions(
	w http.ResponseWriter,
	r *http.Request,
) {
	promotions := []api.PromotionResponse{
		{
			Id:            uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Name:          "Summer Sale",
			DiscountType:  api.PromotionDiscountType("PERCENTAGE"),
			DiscountValue: 20,
			MaxDiscount:   100,
			MaxUses:       1000,
			UsedCount:     150,
			StartsAt:      time.Now(),
			ExpiresAt:     time.Now().Add(30 * 24 * time.Hour),
			Status:        api.PromotionStatus("ACTIVE"),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(promotions)
}

func (h *PromotionHandler) GetPromotionById(
	w http.ResponseWriter,
	r *http.Request,
	id openapi_types.UUID,
) {
	if id != uuid.MustParse("11111111-1111-1111-1111-111111111111") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{
			Error: "Promotion not found",
		})

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(api.PromotionResponse{
		Id:            id,
		Name:          "Summer Sale",
		DiscountType:  "PERCENTAGE",
		DiscountValue: 20,
		MaxDiscount:   100,
		MaxUses:       1000,
		UsedCount:     150,
		StartsAt:      time.Now(),
		ExpiresAt:     time.Now().Add(30 * 24 * time.Hour),
		Status:        "ACTIVE",
	})
}

func (h *PromotionHandler) UpdatePromotion(
	w http.ResponseWriter,
	r *http.Request,
	id openapi_types.UUID,
) {
	var req api.UpdatePromotionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{
			Error: "Invalid request payload",
		})

		return
	}

	if id != uuid.MustParse("11111111-1111-1111-1111-111111111111") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{
			Error: "Promotion not found",
		})

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(api.PromotionResponse{
		Id:            id,
		Name:          req.Name,
		DiscountType:  req.DiscountType,
		DiscountValue: req.DiscountValue,
		MaxDiscount:   req.MaxDiscount,
		MaxUses:       req.MaxUses,
		UsedCount:     150,
		StartsAt:      req.StartsAt,
		ExpiresAt:     req.ExpiresAt,
		Status:        req.Status,
	})
}
