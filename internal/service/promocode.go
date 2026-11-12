package service

import (
	"context"
	"fmt"
	"marketing/internal/dto"

	"github.com/google/uuid"
)

type promoCodeRepo interface {
	CreatePromoCode(
		ctx context.Context,
		promoCodeDTO dto.PromoCode,
	) (
		*uuid.UUID,
		error,
	)
}

type PromoCodeService struct {
	repo promoCodeRepo
}

func NewPromoCodeService(repo promoCodeRepo) *PromoCodeService {
	return &PromoCodeService{
		repo: repo,
	}
}

func (s *PromoCodeService) CreatePromoCode(
	ctx context.Context,
	promoCodeDTO dto.PromoCode,
) (
	*uuid.UUID,
	error,
) {
	id, err := s.repo.CreatePromoCode(ctx, promoCodeDTO)
	if err != nil {
		return nil, fmt.Errorf("service: create promo code: %w", err)
	}

	return id, nil
}
