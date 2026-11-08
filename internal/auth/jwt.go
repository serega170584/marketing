package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

//nolint:gochecknoglobals
var publicKey *ecdsa.PublicKey

//nolint:exhaustruct_v5
func ValidateJWT(
	tokenString string,
	publicKey *ecdsa.PublicKey,
) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodES256 {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}

			return publicKey, nil
		},
		jwt.WithIssuer("marketing-api"),
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

//nolint:gosec
func LoadPublicKey(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return fmt.Errorf("invalid PEM")
	}

	pk, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	lpk, ok := pk.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("not an ECDSA public key")
	}

	publicKey = lpk

	return nil
}

func GetPublicKey() *ecdsa.PublicKey {
	return publicKey
}
