package dto

import "time"

type Promotion struct {
	ID            string                `db:"id"`
	Name          string                `db:"name"`
	DiscountType  PromotionDiscountType `db:"discount_type"`
	DiscountValue int                   `db:"discount_value"`
	MaxDiscount   int                   `db:"max_discount"`
	MaxUses       int64                 `db:"max_uses"`
	UsedCount     int64                 `db:"used_count"`
	StartsAt      time.Time             `db:"starts_at"`
	ExpiresAt     time.Time             `db:"expires_at"`
	Status        PromotionStatus       `db:"status"`
}
