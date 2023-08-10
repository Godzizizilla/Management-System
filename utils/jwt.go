package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var secretKey = []byte("secret_key")

type Claims struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(id uint, role string) (string, int64, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		ID:   id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", 0, err
	}
	return signedToken, claims.IssuedAt, nil
}

func AuthenticateToken(tokenString string) (uint, string, int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// 解析Token错误
	if err != nil {
		return 0, "", 0, errors.New("invalid token")
	}

	// 验证Token是否有效
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.ID, claims.Role, claims.IssuedAt, nil
	}
	return 0, "", 0, errors.New("invalid token")
}
