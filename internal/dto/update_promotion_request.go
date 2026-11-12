package dto

import "time"

type UpdatePromotionRequestDTO struct {
	Name          string                `json:"name"`
	DiscountType  PromotionDiscountType `json:"discount_type"`
	DiscountValue int                   `json:"discount_value"`
	MaxDiscount   int                   `json:"max_discount"`
	MaxUses       int64                 `json:"max_uses"`
	StartsAt      time.Time             `json:"starts_at"`
	ExpiresAt     time.Time             `json:"expires_at"`
	Status        PromotionStatus       `json:"status"`
}
