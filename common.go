package main

import "github.com/dgrijalva/jwt-go"

type JwtClaims struct {
	HasuraClaims HasuraClaims `json:"https://hasura.io/jwt/claim"`
	jwt.StandardClaims
}

type HasuraClaims struct {
	XHasuraUserId string `json:"X-Hasura-User-Id"`
}
