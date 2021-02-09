package client

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

type service struct {
	apiURL string
	auth   string
}

func NewService(apiURL string, user string, password string) (service, error) {
	if apiURL == "" {
		return service{}, fmt.Errorf("mGage apiURL missing")
	}

	// generate basic auth header value from mGage credentials
	var b bytes.Buffer
	input := []byte(fmt.Sprintf("%s:%s", user, password))

	encoder := base64.NewEncoder(base64.StdEncoding, &b)
	_, err := encoder.Write(input)
	if err != nil {
		return service{}, fmt.Errorf("failed encoding: %s", err)
	}
	encoder.Close()

	return service{
		apiURL: apiURL,
		auth:   b.String(),
	}, nil
}
