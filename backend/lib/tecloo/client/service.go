package client

import (
	"fmt"
)

type service struct {
	apiURL string
}

func NewService(apiURL string) (service, error) {
	if apiURL == "" {
		return service{}, fmt.Errorf("tecloo apiURL missing")
	}

	return service{
		apiURL: apiURL,
	}, nil
}
