package dto

import "github.com/google/uuid"

type PromoCode struct {
	ID             string
	Code           string
	PromotionID    uuid.UUID
	MaxUsesPerUser int64
}
