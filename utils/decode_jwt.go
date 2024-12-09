package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Login string
	Role  string
	Exp   int64
	Type  string // как мне сделать штуку типа type TokenType = "access" | "refresh"?
	jwt.MapClaims
}

func DecodeJWT(token string) (Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(
		token,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("1111"), nil
		}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return Claims{}, err
	}

	claims := parsedToken.Claims.(*Claims)

	if claims.Exp < time.Now().Local().Unix() {
		return Claims{}, errors.New("token is expired")
	}

	return Claims{
		Login: claims.Login,
		Role:  claims.Role,
		Exp:   claims.Exp,
		Type:  claims.Type,
	}, nil
}
