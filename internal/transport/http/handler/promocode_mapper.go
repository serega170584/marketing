package handler

import (
	"marketing/internal/dto"
	"marketing/internal/transport/http/api"
)

func mapPromocodeRequest(
	req api.CreatePromoCodeRequest,
) dto.PromoCode {
	return dto.PromoCode{
		Code:           req.Code,
		PromotionID:    req.PromotionId,
		MaxUsesPerUser: req.MaxUsesPerUser,
	}
}
