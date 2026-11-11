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

type PromoRedemptionHandler struct {
	pgConn *database.Postgres
}

func NewPromoRedemptionHandler(
	pgConn *database.Postgres,
) *PromoRedemptionHandler {
	return &PromoRedemptionHandler{
		pgConn: pgConn,
	}
}

func (h *PromoRedemptionHandler) ReservePromoCode(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req api.ReservePromoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{
			Error: "Invalid request payload",
		})

		return
	}

	response := api.PromoRedemptionResponse{
		Id:             uuid.New(),
		PromotionId:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		UserId:         req.UserId,
		OrderId:        req.OrderId,
		DiscountAmount: 2000,
		Status:         api.PromoRedemptionStatus("RESERVED"),
		CreatedAt:      time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *PromoRedemptionHandler) RedeemPromoCode(
	w http.ResponseWriter,
	r *http.Request,
	id openapi_types.UUID,
) {
	if id != uuid.MustParse("33333333-3333-3333-3333-333333333333") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{
			Error: "Promo redemption not found",
		})

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(api.PromoRedemptionResponse{
		Id:             id,
		PromotionId:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		UserId:         1,
		OrderId:        42,
		DiscountAmount: 2000,
		Status:         "REDEEMED",
		CreatedAt:      time.Now(),
	})
}

func (h *PromoRedemptionHandler) ReleasePromoCode(
	w http.ResponseWriter,
	r *http.Request,
	id openapi_types.UUID,
) {
	if id != uuid.MustParse("33333333-3333-3333-3333-333333333333") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{
			Error: "Promo redemption not found",
		})

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(api.PromoRedemptionResponse{
		Id:             id,
		PromotionId:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		UserId:         1,
		OrderId:        42,
		DiscountAmount: 2000,
		Status:         "RELEASED",
		CreatedAt:      time.Now(),
	})
}

func (h *PromoRedemptionHandler) GetPromoRedemptions(
	w http.ResponseWriter,
	r *http.Request,
	params api.GetPromoRedemptionsParams,
) {
	redemptions := []api.PromoRedemptionResponse{
		{
			Id:             uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			PromotionId:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			UserId:         1,
			OrderId:        42,
			DiscountAmount: 2000,
			Status:         api.PromoRedemptionStatus("REDEEMED"),
			CreatedAt:      time.Now(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(redemptions)
}
