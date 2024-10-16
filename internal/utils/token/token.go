package token

import (
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(userID uint, expiresAt time.Time) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(), // token expire time
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("default")) // key to type []byte
	if err != nil {
		zap.S().Errorf("Error signing token: %v", err)
		return "", err
	}

	return signedToken, nil
}
