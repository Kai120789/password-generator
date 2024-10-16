package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

// Claims для JWT
type Claims struct {
	UserID uint `json:"id"`
	jwt.StandardClaims
}

// middleware for Access token check
func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract token from header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "token is missing", http.StatusUnauthorized)
			return
		}

		// check is token start with 'Bearer '
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "invalid token format", http.StatusUnauthorized)
			return
		}

		// trim prefix 'Bearer '
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// where save info from token
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// check method is true
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("method is not correct: %v", token.Header["alg"])
			}
			return []byte("default"), nil
		})

		// check token is valid
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		zap.S().Infof("claims: %+v", claims)

		ctx := context.WithValue(r.Context(), "id", claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
