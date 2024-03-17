package middleware

import (
	"auth_service/pkg/utils"
	"context"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Authorization header must be in te format `Bearer {token}`", http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]
		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userData", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
