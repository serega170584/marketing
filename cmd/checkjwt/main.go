package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

//nolint:exhaustruct_v5
func main() {
	var rootCmd = &cobra.Command{
		Use:   "jwt",
		Short: "JWT checker",
	}

	var checkJWTCmd = &cobra.Command{
		Use:   "check",
		Short: "check JWT",
		Run: func(cmd *cobra.Command, args []string) {
			publicKey, err := LoadPublicKey("public.pem")
			if err != nil {
				fmt.Println(err)

				os.Exit(1)
			}

			tokenStr, err := LoadToken("token")
			if err != nil {
				fmt.Println(err)

				os.Exit(1)
			}

			token, err := ValidateJWT(tokenStr, publicKey)
			if err != nil {
				fmt.Println(err)

				os.Exit(1)
			}

			if !token.Valid {
				fmt.Println("invalid token")

				os.Exit(1)
			}

			fmt.Println("valid token")
		},
	}

	rootCmd.AddCommand(checkJWTCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func ValidateJWT(
	tokenString string,
	publicKey *ecdsa.PublicKey,
) (*jwt.Token, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodES256 {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}

			return publicKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

//nolint:gosec
func LoadPublicKey(path string) (
	*ecdsa.PublicKey,
	error,
) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}

	return key, nil
}

//nolint:gosec
func LoadToken(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	token := strings.TrimSpace(string(data))

	if token == "" {
		return "", fmt.Errorf("token file is empty")
	}

	return token, nil
}
