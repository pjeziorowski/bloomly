package api

import (
	"bloomly/backend/config"
	"github.com/hasura/go-graphql-client"
	"net/http"
)

var HasuraClient = graphql.NewClient(config.HasuraApiUrl, &http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport}})

type transport struct {
	underlyingTransport http.RoundTripper
}

func (t transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("x-hasura-admin-secret", config.HasuraApiToken)
	req.Header.Add("content-type", "application/json")
	return t.underlyingTransport.RoundTrip(req)
}
