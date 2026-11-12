package auth

import (
	"context"
	"slices"
)

import "github.com/golang-jwt/jwt/v5"

const PromoCodeCreatePermission = "promo_code:create"
const PromoCodeViewPermission = "promo_code:view"
const PromotionCreatePermission = "promotion:create"
const PromotionViewPermission = "promotion:view"
const PromoRedemptionCreatePermission = "promo_redemption:create"
const PromoRedemptionViewPermission = "promo_redemption:view"

type Claims struct {
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type ClaimsContextKey struct{}

func (c Claims) HasPermission(permission string) bool {
	return slices.Contains(c.Permissions, permission)
}

func GetClaims(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey{}).(*Claims)
	return claims, ok
}
