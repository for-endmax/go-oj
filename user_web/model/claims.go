package model

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID       uint
	Nickname string
	Role     int32
	jwt.StandardClaims
}
