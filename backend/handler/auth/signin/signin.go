package signin

import (
	"bloomly/backend/api"
	"bloomly/backend/handler/auth"
	"context"
	"encoding/json"
	"errors"
	"github.com/hasura/go-graphql-client"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
)

type ActionPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            SigninArgs             `json:"input"`
}

type GraphQLError struct {
	Message string `json:"message"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// set the response header as JSON
	w.Header().Set("Content-Type", "application/json")

	// read request body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// parse the body as action payload
	var actionPayload ActionPayload
	err = json.Unmarshal(reqBody, &actionPayload)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// Send the request params to the Action's generated handler function
	result, err := signin(actionPayload.Input)

	// throw if an error happens
	if err != nil {
		errorObject := GraphQLError{
			Message: err.Error(),
		}
		errorBody, _ := json.Marshal(errorObject)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorBody)
		return
	}

	// Write the response as JSON
	data, _ := json.Marshal(result)
	w.Write(data)
}

func signin(args SigninArgs) (response SigninOutput, err error) {
	log.Printf("received sign in request %v", args)

	response = SigninOutput{}

	// checking user password
	id, err := verify(args.Input.Email, args.Input.Password)
	if err != nil {
		return response, err
	}

	// try to create and sign a token
	signed, err := auth.SignTokenFor(id)
	if err != nil {
		return response, err
	}
	response.AccessToken = signed

	return response, nil
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

	// getting user by email
	err := api.HasuraClient.Query(context.Background(), &query, vars)
	if err != nil {
		return 0, err
	}
	if len(query.Creators) != 1 {
		return 0, errors.New("user with email " + email + " not found")
	}

	// checking password, returning user ID if fine
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
