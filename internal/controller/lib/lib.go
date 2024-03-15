package lib

import (
	"crypto/sha1"
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
