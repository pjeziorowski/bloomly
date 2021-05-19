package config

import (
	"errors"
	"os"
)

var (
	Secret         = os.Getenv("SECRET")
	HasuraApiUrl   = os.Getenv("HASURA_API_URL")
	HasuraApiToken = os.Getenv("HASURA_API_TOKEN")
)

func CheckServerConfig() []error {
	var configErrors []error

	if Secret == "" {
		configErrors = append(configErrors, errors.New("SECRET env required"))
	}
	if HasuraApiUrl == "" {
		configErrors = append(configErrors, errors.New("HASURA_API_URL env required"))
	}
	if HasuraApiToken == "" {
		configErrors = append(configErrors, errors.New("HASURA_API_TOKEN env required"))
	}

	return configErrors
}
