package common

import "github.com/golang-jwt/jwt/v5"

type ctxKey string

const (
	JwtCtxKey            ctxKey = "jwtContextKey"
	EncodedUserJwtCtxKey ctxKey = "encodedUserJwtCtxKey"
)

func (c ctxKey) ToString() string {
	return string(c)
}

type UserClaims struct {
	Id int64 `json:"id"`
	jwt.RegisteredClaims
}