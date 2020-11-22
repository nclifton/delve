package client

import (
	"fmt"
	"net/http"
)

type Service struct {
	apiURL string
	http   *http.Client
}

func NewService(apiURL string, client *http.Client) (*Service, error) {
	if apiURL == "" {
		return &Service{}, fmt.Errorf("alaris apiURL missing")
	}

	return &Service{
		apiURL: apiURL,
		http:   client,
	}, nil
}
