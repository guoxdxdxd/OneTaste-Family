package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"onetaste-family/backend/internal/config"
)

var (
	// ErrInvalidToken token无效
	ErrInvalidToken = errors.New("invalid token")
)

// Claims JWT Claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(config.AppConfig.JWT.Expiration)
	
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetTokenExpiration 获取Token过期时间（秒）
func GetTokenExpiration() int64 {
	return int64(config.AppConfig.JWT.Expiration.Seconds())
}
