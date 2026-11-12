package repository

import (
	"context"
	"fmt"
	"marketing/internal/dto"
	"marketing/internal/infrastructure/database"

	"github.com/google/uuid"
)

type PromoCodeRepo struct {
	postgres *database.Postgres
}

func NewPromoCodeRepo(
	postgres *database.Postgres,
) *PromoCodeRepo {
	return &PromoCodeRepo{
		postgres: postgres,
	}
}

func (r PromoCodeRepo) CreatePromoCode(
	ctx context.Context,
	promoCodeDTO dto.PromoCode,
) (*uuid.UUID, error) {
	pool := r.postgres.Pool()

	id := uuid.New()

	query := `
        INSERT INTO promo_codes (
            id,
            code,
            promotion_id,
            max_uses_per_user
        )
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	err := pool.QueryRow(
		ctx,
		query,
		id,
		promoCodeDTO.Code,
		promoCodeDTO.PromotionID,
		promoCodeDTO.MaxUsesPerUser,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create promo code: %w", err)
	}

	return &id, nil
}
