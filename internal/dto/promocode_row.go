package dto

import "github.com/google/uuid"

type PromoCodeRow struct {
	ID             string    `db:"id"`
	Code           string    `db:"code"`
	PromotionID    uuid.UUID `db:"promotion_id"`
	MaxUsesPerUser int64     `db:"max_uses_per_user"`
}
