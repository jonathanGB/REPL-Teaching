package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"os"
)

var JWT_SECRET []byte = []byte(os.Getenv("JWT_SECRET"))

func MarshalToken(name, id, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
		"name": name,
		"role": role,
		"expire": time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString(JWT_SECRET)
}
