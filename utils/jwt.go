package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var SecretKey []byte

type Claims struct {
	UserID string `json:"id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(userID string, role string) (token string, jti string, err error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	jtiStr, err := generateRandomString(16)
	if err != nil {
		return "", "", err
	}
	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        jtiStr,
		},
	}

	genToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := genToken.SignedString(SecretKey)
	if err != nil {
		return "", "", err
	}
	return signedToken, claims.Id, nil
}

func AuthenticateToken(tokenString string) (userID string, role string, jti string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	// 解析Token错误
	if err != nil {
		return "", "", "", errors.New("invalid token")
	}

	// 验证Token是否有效
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, claims.Role, claims.Id, nil
	}
	return "", "", "", errors.New("invalid token")
}

// generateRandomString generates a random string of the given length.
func generateRandomString(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
