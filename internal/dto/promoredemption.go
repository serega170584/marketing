package dto

import "time"

type PromoRedemption struct {
	ID             string                `db:"id"`
	PromotionID    string                `db:"promotion_id"`
	UserID         string                `db:"user_id"`
	OrderID        string                `db:"order_id"`
	DiscountAmount int64                 `db:"discount_amount"`
	Status         PromoRedemptionStatus `db:"status"`
	CreatedAt      time.Time             `db:"created_at"`
}
