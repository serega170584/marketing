package dto

type CreatePromoCodeRequest struct {
	Code           string `json:"code"`
	PromotionID    string `json:"promotion_id"`
	MaxUsesPerUser int64  `json:"max_uses_per_user"`
}
