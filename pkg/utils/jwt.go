package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var mySigningKey = []byte("rbac_system_secret")

type Myclaim struct {
	UserId int64 `json:"user_id`
	jwt.RegisteredClaims
}

func GenerateToken(userId int64) (string, error) {

	c := Myclaim{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "rbac_system",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := token.SignedString(mySigningKey)
	return ss, err
}
