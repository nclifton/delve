package client

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
)

type service struct {
	apiURL string
	auth   string
}

func NewService(apiURL string, optusUser string, optusPass string) (service, error) {
	if apiURL == "" {
		return service{}, fmt.Errorf("optus apiURL missing")
	}

	optusUser = optusUser
	optusPass = optusPass

	// generate basic auth header value from optus credentials
	var b bytes.Buffer
	input := []byte(fmt.Sprintf("%s:%s", optusUser, optusPass))

	encoder := base64.NewEncoder(base64.StdEncoding, &b)
	_, err := encoder.Write(input)
	if err != nil {
		log.Fatalln("failed encoding:", err)
	}
	encoder.Close()

	return service{
		apiURL: apiURL,
		auth:   b.String(),
	}, nil
}
