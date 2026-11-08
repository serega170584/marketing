package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

//nolint:exhaustruct_v5
func main() {
	var rootCmd = &cobra.Command{
		Use:   "jwt",
		Short: "JWT generator",
	}

	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate JWT",
		Run: func(cmd *cobra.Command, args []string) {
			privateKey, err := LoadPrivateKey("private.pem")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			token, err := GenerateJWT(privateKey)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println(token)
		},
	}

	rootCmd.AddCommand(generateCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

//nolint:gosec
func LoadPrivateKey(path string) (
	*ecdsa.PrivateKey,
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

	return x509.ParseECPrivateKey(block.Bytes)
}

func GenerateJWT(privateKey *ecdsa.PrivateKey) (
	string,
	error,
) {
	claims := jwt.MapClaims{
		"iss": "marketing-api",
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"permissions": []string{
			"users:read",
			"users:create",
			"users:update",
			"orders:read",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString(privateKey)
}
