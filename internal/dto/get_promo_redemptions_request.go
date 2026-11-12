package dto

type GetPromoRedemptionsRequest struct {
	UserID      *string                `json:"user_id"`
	PromotionID *string                `json:"promotion_id"`
	OrderID     *string                `json:"order_id"`
	Status      *PromoRedemptionStatus `json:"status"`
}
