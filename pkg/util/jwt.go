package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte("SUPER_SECRET_KEY") // 🔒 move to ENV later

type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokens(userID int64, role string) (string, string, error) {
	// Access token (15m)
	accessClaims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(JWTSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh token (7d)
	refreshClaims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(JWTSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
