package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SendMMSParams struct {
	AllowAdaptation string    `json:"allowAdaptation"`
	CustomerID      string    `json:"customerId"`
	ForwardLock     string    `json:"forwardLock"`
	From            string    `json:"from"`
	Messages        []Message `json:"message"`
	RefMessageID    string    `json:"refMessageId"`
	ReportingKey1   string    `json:"reportingKey1"`
	ReportingKey2   string    `json:"reportingKey2"`
	RequestID       string    `json:"requestId"`
	Subject         string    `json:"subject"`
	To              []string  `json:"to"`
}

type Message struct {
	Text string `json:"text"`

	DisplayName string `json:"displayName"`
	MediaType   string `json:"mediaType"`
	MediaURL    string `json:"mediaUrl"`

	TranscodeId string `json:"transcodeId"`
}

type SendMMSResponse struct {
	RecipientDetails []RecipientDetails `json:"recipientDetails"`
	RequestID        string             `json:"requestId"`
}

type SendMMSErrorResponse struct {
	RecipientDetails []RecipientDetails `json:"recipientDetails"`
	RequestID        string             `json:"requestId"`
	Status           string             `json:"status"`
	ErrorText        string             `json:"errorText"`
}

type RecipientDetails struct {
	MessageID string `json:"messageId"`
	To        string `json:"to"`
}

func (s service) SendMMS(ctx context.Context, params SendMMSParams) (SendMMSResponse, int, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return SendMMSResponse{}, 0, err
	}

	req, err := http.NewRequest("POST", s.apiURL+"/wsgw/v1/mms/messages", bytes.NewBuffer(b))
	if err != nil {
		return SendMMSResponse{}, 0, err
	}

	req.Header.Add("Authorization", "Basic "+s.auth)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return SendMMSResponse{}, 0, err
	}

	if res.StatusCode != 200 {
		output := SendMMSErrorResponse{}
		if err := json.NewDecoder(res.Body).Decode(&output); err != nil {
			return SendMMSResponse{}, res.StatusCode, err
		}

		return SendMMSResponse{}, res.StatusCode, fmt.Errorf("mGage gave status code %d, response: %s", res.StatusCode, output.ErrorText)
	}

	output := SendMMSResponse{}
	if err := json.NewDecoder(res.Body).Decode(&output); err != nil {
		return SendMMSResponse{}, 0, err
	}

	return output, res.StatusCode, nil
}
