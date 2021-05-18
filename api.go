package main

import (
	"github.com/hasura/go-graphql-client"
	"net/http"
)

type Transport struct {
	UnderlyingTransport http.RoundTripper
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("x-hasura-admin-secret", hasuraApiToken)
	req.Header.Add("content-type", "application/json")
	return t.UnderlyingTransport.RoundTrip(req)
}

var client = graphql.NewClient(hasuraApiUrl, &http.Client{Transport: &Transport{UnderlyingTransport: http.DefaultTransport}})
