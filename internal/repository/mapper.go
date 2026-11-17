package repository

import (
	"marketing/internal/dto"
)

func mapPromoCode(
	promoCode dto.PromoCode,
) dto.PromoCodeRow {
	return dto.PromoCodeRow{
		ID:             promoCode.ID,
		Code:           promoCode.Code,
		PromotionID:    promoCode.PromotionID,
		MaxUsesPerUser: promoCode.MaxUsesPerUser,
	}
}
