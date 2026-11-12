package dto

type PromotionDiscountType string

const (
	PromotionDiscountTypePercentage PromotionDiscountType = "PERCENTAGE"
	PromotionDiscountTypeFixed      PromotionDiscountType = "FIXED"
)

type PromotionStatus string

const (
	PromotionStatusActive   PromotionStatus = "ACTIVE"
	PromotionStatusDisabled PromotionStatus = "DISABLED"
)

type PromoRedemptionStatus string

const (
	PromoRedemptionStatusReserved PromoRedemptionStatus = "RESERVED"
	PromoRedemptionStatusRedeemed PromoRedemptionStatus = "REDEEMED"
	PromoRedemptionStatusReleased PromoRedemptionStatus = "RELEASED"
)
