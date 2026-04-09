package token

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var secret = []byte("CHANGE_ME")

func Generate(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}
