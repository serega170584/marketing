package usecase

import (
	"context"
	"marketing/internal/dto"

	"github.com/google/uuid"
)

type promoCodeService interface {
	CreatePromoCode(
		ctx context.Context,
		codeDTO dto.PromoCode,
	) (
		*uuid.UUID,
		error,
	)
}

type PromoCodeUseCase struct {
	promoCodeService promoCodeService
}

func NewPromoCodeUseCase(
	promoCodeService promoCodeService,
) *PromoCodeUseCase {
	return &PromoCodeUseCase{
		promoCodeService: promoCodeService,
	}
}

func (uc *PromoCodeUseCase) CreatePromoCode(
	ctx context.Context,
	promoCodeDTO dto.PromoCode,
) (
	*uuid.UUID,
	error,
) {
	id, err := uc.promoCodeService.CreatePromoCode(
		ctx,
		promoCodeDTO,
	)
	if err != nil {
		return nil, err
	}

	return id, nil
}
