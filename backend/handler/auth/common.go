package auth

import (
	"bloomly/backend/config"
	"bloomly/backend/handler"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hasura/go-graphql-client"
	"time"
)

func SignTokenFor(userId graphql.Int) (string, error) {
	claims := &handler.JwtClaims{
		handler.HasuraClaims{XHasuraUserId: fmt.Sprintf("%v", userId)},
		jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 72).Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.Secret))

	if err != nil {
		return "", err
	}

	return signed, nil
}
