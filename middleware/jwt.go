package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error": "Unauthorized", "message": "Token JWT ausente ou invalido"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			parser := jwt.NewParser()
			
			// Decodifica o token sem exigir a verificação da assinatura criptográfica
			token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
			if err != nil {
				http.Error(w, `{"error": "Unauthorized", "message": "Token JWT invalido ou malformatado"}`, http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if userID, exists := claims["sub"]; exists {
					r.Header.Set("X-User-Id", fmt.Sprintf("%v", userID))
				}
				if role, exists := claims["role"]; exists {
					r.Header.Set("X-User-Role", fmt.Sprintf("%v", role))
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}