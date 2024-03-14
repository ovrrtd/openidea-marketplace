package jwt

import (
	"ecomm/internal/helper/errorer"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwt(payload jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte("sampleSecretKey"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyJwt(tokenString string, claims jwt.Claims) error {
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte("sampleSecretKey"), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return errorer.ErrUnauthorized
		}
		return err

	}
	if !tkn.Valid {
		return errorer.ErrUnauthorized
	}

	return nil
}
