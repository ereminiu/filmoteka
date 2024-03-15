package lib

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const salt = "I192020LOVE19248K@tya38302Alekseeva1997"
const signedString = "I18203048LOVE32849K@tya38349Alekseeva1997"

func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func GenerateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(12 * time.Hour)},
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
	},
		userId,
	})

	return token.SignedString([]byte(signedString))
}

func ParseToken(inputToken string) (int, error) {
	token, err := jwt.ParseWithClaims(inputToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signedString), nil
	})
	if err != nil {
		return -1, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return -1, errors.New("invalid token claims type")
	}
	return claims.UserId, nil
}
