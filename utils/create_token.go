package utils

import "github.com/golang-jwt/jwt/v5"

func CreateToken(exp int64, secret string, login string, role string, token_type string) (string, error) {
	access := jwt.New(jwt.SigningMethodHS256)
	claims := access.Claims.(jwt.MapClaims)
	claims["Login"] = login
	claims["Role"] = role
	claims["Exp"] = exp
	claims["Type"] = token_type
	tokenString, err := access.SignedString([]byte("1111"))
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}
