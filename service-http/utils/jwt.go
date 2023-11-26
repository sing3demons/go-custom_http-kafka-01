package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getSecretKeyFromEnv() (privateKey []byte, publicKey []byte, err error) {
	private := os.Getenv("PRIVATE_KEY")
	if private == "" {
		return nil, nil, fmt.Errorf("private key not found")
	}
	privateKey, err = base64.StdEncoding.DecodeString(private)
	if err != nil {
		return nil, nil, err
	}

	public := os.Getenv("PUBLIC_KEY")
	if public == "" {
		return nil, nil, fmt.Errorf("public key not found")
	}
	publicKey, err = base64.StdEncoding.DecodeString(public)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, publicKey, nil
}

func GenerateToken(sub string) (token string, err error) {
	issuer := os.Getenv("JWT_ISSUER")
	var audience jwt.ClaimStrings
	aud := os.Getenv("AUDIENCE")
	if aud != "" {
		audience = strings.Split(aud, ",")
	}

	privateKey, _, err := getSecretKeyFromEnv()
	if err != nil {
		return "", err
	}

	rsa, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	claims := jwt.RegisteredClaims{
		Subject:   sub,
		Issuer:    issuer,
		Audience:  audience,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(rsa)
}
