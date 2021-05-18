package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hasura/go-graphql-client"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Signin(c echo.Context) error {
	input := new(SigninArgs)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	c.Logger().Print("Received sign in request %v", input)

	id, err := verify(input.Arg1.Email, input.Arg1.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	claims := &JwtClaims{
		HasuraClaims{XHasuraUserId: fmt.Sprintf("%v", id)},
		jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 72).Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	output := new(SigninOutput)
	output.AccessToken = signed

	return c.JSON(http.StatusOK, output)
}

func verify(email, password string) (graphql.Int, error) {
	var query struct {
		Creators []struct {
			Id       graphql.Int
			Password graphql.String
		} `graphql:"creators(where: {verified_email: {_eq: $verified_email}})"`
	}
	vars := map[string]interface{}{
		"verified_email": graphql.String(email),
	}

	err := client.Query(context.Background(), &query, vars)
	if err != nil {
		return 0, err
	}
	if len(query.Creators) != 1 {
		return 0, errors.New("User with email " + email + " not found. ")
	}

	match := compare([]byte(query.Creators[0].Password), []byte(password))
	if match {
		return query.Creators[0].Id, nil
	}

	return 0, errors.New("Password doesn't match. ")
}

func compare(hashed []byte, plain []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashed, plain)
	if err != nil {
		return false
	}
	return true
}
