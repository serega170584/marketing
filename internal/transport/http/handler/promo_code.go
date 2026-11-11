package handler

import (
	"encoding/json"
	"net/http"

	"marketing/internal/infrastructure/database"
	"marketing/internal/transport/http/api"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type PromoCodeHandler struct {
	pgConn *database.Postgres
}

func NewPromoCodeHandler(pgConn *database.Postgres) *PromoCodeHandler {
	return &PromoCodeHandler{
		pgConn: pgConn,
	}
}

func (h *PromoCodeHandler) CreatePromoCode(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req api.CreatePromoCodeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{
			Error: "Invalid request payload",
		})

		return
	}

	response := api.PromoCodeResponse{
		Id:             uuid.New(),
		Code:           req.Code,
		PromotionId:    req.PromotionId,
		MaxUsesPerUser: req.MaxUsesPerUser,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *PromoCodeHandler) GetPromoCodes(
	w http.ResponseWriter,
	r *http.Request,
	params api.GetPromoCodesParams,
) {
	codes := []api.PromoCodeResponse{
		{
			Id:             uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			Code:           "SUMMER2026",
			PromotionId:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			MaxUsesPerUser: 1,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(codes)
}

func (h *PromoCodeHandler) GetPromoCodeById(
	w http.ResponseWriter,
	r *http.Request,
	id openapi_types.UUID,
) {
	if id != uuid.MustParse("22222222-2222-2222-2222-222222222222") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{
			Error: "Promo code not found",
		})

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(api.PromoCodeResponse{
		Id:             id,
		Code:           "SUMMER2026",
		PromotionId:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		MaxUsesPerUser: 1,
	})
}
