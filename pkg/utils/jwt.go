package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var mySigningKey = []byte("rbac_system_secret")

const TokenTTL = 2 * time.Hour

type Myclaim struct {
	UserId int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userId int64) (string, error) {

	c := Myclaim{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "rbac_system",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := token.SignedString(mySigningKey)
	return ss, err
}

func ParseToken(tokenString string) (int64, time.Time, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Myclaim{}, func(t *jwt.Token) (any, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return 0, time.Time{}, err
	} else if claims, ok := token.Claims.(*Myclaim); ok && token.Valid && claims.ExpiresAt != nil {
		return claims.UserId, claims.ExpiresAt.Time, nil
	}
	return 0, time.Time{}, errors.New("invalid token")
}
