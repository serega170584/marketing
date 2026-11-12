package dto

type PromoCode struct {
	ID             string `db:"id"`
	Code           string `db:"code"`
	PromotionID    string `db:"promotion_id"`
	MaxUsesPerUser int64  `db:"max_uses_per_user"`
}
