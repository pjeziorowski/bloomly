package signup

import (
	"bloomly/backend/api"
	"bloomly/backend/handler/auth"
	"context"
	"encoding/json"
	"github.com/hasura/go-graphql-client"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
)

type ActionPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            SignupArgs             `json:"input"`
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
	result, err := signup(actionPayload.Input)

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

func signup(args SignupArgs) (response SignupOutput, err error) {
	log.Printf("received sign up request %v", args)

	response = SignupOutput{}

	// try to create a new user
	id, err := register(args.Input.Email, args.Input.Password)
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

func register(email, password string) (graphql.Int, error) {
	var mutation struct {
		InsertCreatorsOne struct {
			Id graphql.Int
		} `graphql:"insert_creators_one(object: {verified_email: $verified_email, password: $password})"`
	}
	// encrypting password not to store it in plain text
	hashed, err := hash([]byte(password))
	if err != nil {
		return 0, err
	}
	vars := map[string]interface{}{
		"password":       graphql.String(hashed),
		"verified_email": graphql.String(email),
	}

	// trying to create a new user in Hasura backend
	err = api.HasuraClient.Mutate(context.Background(), &mutation, vars)
	if err != nil {
		return 0, err
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
