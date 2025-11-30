package service

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func validateToken(tokenString, jwtSecret string) (int, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return 0, "", fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, "", fmt.Errorf("invalid token claims")
		}
		username, ok := claims["username"].(string)
		if !ok {
			return 0, "", fmt.Errorf("invalid token claims")
		}
		return int(userID), username, nil
	}

	return 0, "", fmt.Errorf("invalid token")
}

