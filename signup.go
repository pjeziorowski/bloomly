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

func Signup(c echo.Context) error {
	input := new(SignupInput)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := addCreator(input.Email, input.Password)
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

	output := new(SignupOutput)
	output.AccessToken = signed

	return c.JSON(http.StatusOK, output)
}

func addCreator(email, password string) (graphql.Int, error) {
	var mutation struct {
		InsertCreatorsOne struct {
			Id graphql.Int
		} `graphql:"insert_creators_one(object: {verified_email: $verified_email, password: $password})"`
	}
	hashed, err := hash([]byte(password))
	if err != nil {
		return 0, err
	}
	vars := map[string]interface{}{
		"password":       graphql.String(hashed),
		"verified_email": graphql.String(email),
	}

	err = client.Mutate(context.Background(), &mutation, vars)
	if err != nil {
		return 0, errors.New("Error creating a new account. ")
	}

	return mutation.InsertCreatorsOne.Id, nil
}

func hash(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
