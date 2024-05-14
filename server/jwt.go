package server

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey          = []byte(os.Getenv("JWT_SECRET_KEY"))
	ErrInvalidJwtToken = errors.New("invalid jwt token")
)

const jwtTokenDuration = 10 * time.Minute

func generateJWT(username string) (map[string]any, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":   username,
		"exp":        time.Now().Add(jwtTokenDuration),
		"authorized": true,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return map[string]any{}, fmt.Errorf("error generating jwt token: %w", err)
	}

	return map[string]any{
		"access_token": tokenString,
		"token_type":   "Bearer",
		"expires_in":   jwtTokenDuration,
	}, nil
}

func parseAndValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return jwt.MapClaims{}, fmt.Errorf("error parsing toknenString during validation: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return jwt.MapClaims{}, ErrInvalidJwtToken
	}
}
