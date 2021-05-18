package main

import (
	"errors"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
)

var (
	secret         = os.Getenv("SECRET")
	hasuraApiUrl   = os.Getenv("HASURA_API_URL")
	hasuraApiToken = os.Getenv("HASURA_API_TOKEN")
)

func checkEnv() []error {
	var configErrors []error

	if secret == "" {
		configErrors = append(configErrors, errors.New("SECRET env required. "))
	}
	if hasuraApiUrl == "" {
		configErrors = append(configErrors, errors.New("HASURA_API_URL env required. "))
	}
	if hasuraApiToken == "" {
		configErrors = append(configErrors, errors.New("HASURA_API_TOKEN env required. "))
	}

	return configErrors
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	configErrors := checkEnv()
	if configErrors != nil {
		for _, err := range configErrors {
			e.Logger.Error(err.Error())
		}
		e.Logger.Fatal("Killing the server. ")
		return
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/signin", Signin)
	e.POST("/signup", Signup)

	e.Logger.Fatal(e.Start(":1323"))
}
