package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JwtSecret []byte // 将从配置读取

// SetJwtSecret 设置 JWT Secret
func SetJwtSecret(secret string) {
	JwtSecret = []byte(secret)
}

// CustomClaims 自定义JWT声明
type CustomClaims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(id int64, username, phone string, expirationHours int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)
	claims := &CustomClaims{
		ID:       id,
		Username: username,
		Phone:    phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "power-admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
