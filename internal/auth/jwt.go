package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtExpiry = 24 * time.Hour

type Claims struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	IsAdmin   bool   `json:"is_admin"`
	AvatarURL string `json:"avatar_url"`
	jwt.RegisteredClaims
}

// Sign creates a signed JWT for the given user fields.
func Sign(userID uint, username string, isAdmin bool, avatarURL string, secret string) (string, error) {
	claims := Claims{
		UserID:    userID,
		Username:  username,
		IsAdmin:   isAdmin,
		AvatarURL: avatarURL,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Verify parses and validates a JWT, returning its claims.
func Verify(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
