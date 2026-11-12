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
) (
	*uuid.UUID,
	error,
) {
	pool := r.postgres.Pool()

	query := `
        UPDATE promo_codes
        SET
            code = $1,
            promotion_id = $2,
            max_uses_per_user = $3
        WHERE id = $4
        RETURNING id
    `

	var id uuid.UUID

	err := pool.QueryRow(
		ctx,
		query,
		promoCodeDTO.Code,
		promoCodeDTO.PromotionID,
		promoCodeDTO.MaxUsesPerUser,
		promoCodeDTO.ID,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to update promo code: %w", err)
	}

	return &id, nil
}
