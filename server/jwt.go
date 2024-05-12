package server

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

const jwtTokenDuration = 10 * time.Minute

func generateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":   username,
		"exp":        time.Now().Add(jwtTokenDuration),
		"authorized": true,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error generating jwt token: %w", err)
	}

	return tokenString, nil
}

func parseJWT(tokenString string) (jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return jwt.Token{}, fmt.Errorf("error parsing toknenString during validation: %w", err)
	}

	return *token, nil
}
