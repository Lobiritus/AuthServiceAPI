package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

// Claims - структура для JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWT генерирует новый JWT токен для пользователя

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour) // Токен истекает через час
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	jwtKey := os.Getenv("JWT_SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))

	return tokenString, err
}

// Проверка токена
func VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	jwtKey := os.Getenv("JWT_SECRET_KEY")
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorMalformed)
		}
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.NewValidationError("invalid token", jwt.ValidationErrorSignatureInvalid)
	}

	// Проверяем, не истек ли срок действия токена
	if !claims.ExpiresAt.Time.After(time.Now()) {
		return nil, jwt.NewValidationError("expired token", jwt.ValidationErrorExpired)
	}

	return claims, nil
}
