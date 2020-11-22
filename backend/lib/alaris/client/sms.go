package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type alarisResponse struct {
	MessageID  string `json:"message_id"`
	DNIS       string `json:"dnis"`
	SegmentNum string `json:"segment_num"`
}

type SendSMSParams struct {
	Username        string
	Password        string
	Command         string
	Message         string
	DNIS            string
	ANI             string
	LongMessageMode string
}

func (s *Service) SendSMS(params SendSMSParams) (string, error) {
	// hand off to alaris
	req, err := http.NewRequest("GET", s.apiURL, nil)
	if err != nil {
		return "", NewError(err.Error(), false)
	}

	q := url.Values{}
	q.Add("username", params.Username)
	q.Add("password", params.Password)
	q.Add("command", params.Command)
	q.Add("message", params.Message)
	q.Add("dnis", params.DNIS)
	if params.ANI != "" {
		q.Add("ani", params.ANI)
	}
	q.Add("longMessageMode", params.LongMessageMode)
	req.URL.RawQuery = q.Encode()

	resp, err := s.http.Do(req)
	if err != nil {
		return "", NewError(err.Error(), true)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", NewError(err.Error(), true)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			if string(respBody) == "NO ROUTES" {
				return "", NewError("NO ROUTES returned from alaris", false)
			}
		}
		return "", NewError(fmt.Sprintf("Failed sending sms, alaris responded with status: %d response: %+v", resp.StatusCode, string(respBody)), false)
	}

	var messageID string

	if string(respBody[0]) == "{" {
		alarisResponse := alarisResponse{}
		err = json.NewDecoder(bytes.NewReader(respBody)).Decode(&alarisResponse)
		if err != nil {
			return "", NewError(err.Error(), true)
		}

		messageID = alarisResponse.MessageID
	}

	if string(respBody[0]) == "[" {
		log.Printf("Response: %s", string(respBody))
		alarisResponse := []alarisResponse{}
		err = json.NewDecoder(bytes.NewReader(respBody)).Decode(&alarisResponse)
		if err != nil {
			return "", NewError(err.Error(), true)
		}

		if len(alarisResponse) > 0 {
			messageID = alarisResponse[0].MessageID
		}
	}

	return messageID, nil
}
