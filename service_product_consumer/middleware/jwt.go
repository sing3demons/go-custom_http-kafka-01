package middleware

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	public := os.Getenv("PUBLIC_KEY")
	if public == "" {
		return nil, fmt.Errorf("public key not found")
	}
	publicKey, err := base64.StdEncoding.DecodeString(public)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return jwt.ParseRSAPublicKeyFromPEM(publicKey)
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	iss, err := claims.GetIssuer()
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	expectedIssuer := os.Getenv("JWT_ISSUER")
	if iss != expectedIssuer {
		return nil, fmt.Errorf("validate: invalid issuer")
	}

	expectedAudience := os.Getenv("AUDIENCE")
	audience, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	foundAudience := false
	for _, aud := range audience {
		if aud == expectedAudience {
			foundAudience = true
			break
		}
	}

	if !foundAudience {
		return nil, fmt.Errorf("invalid audience. Expected: %s, Got: %v", expectedAudience, audience)
	}

	return claims, nil
}
