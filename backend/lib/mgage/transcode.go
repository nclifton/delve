package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

type CreateTranscodesWithURLParams struct {
	CustomerID      string `json:"customerId"`
	RequestID       string `json:"requestId"`
	MediaType       string `json:"mediaType"`
	MediaURL        string `json:"mediaUrl"`
	NotificationURL string `json:"notificationUrl"`
}

type CreateTranscodesWithContentParams struct {
	CustomerID      string
	RequestID       string
	MediaType       string
	NotificationURL string
	ContentType     string
	FileName        string
	Content         []byte
}

type CreateTranscodesResponse struct {
	RequestID   string `json:"requestId"`
	TranscodeID string `json:"transcodeId"`
}

type CreateTranscodesErrorResponse struct {
	RequestID string `json:"requestId"`
	Status    string `json:"status"`
	ErrorText string `json:"errorText"`
}

func (s service) CreateTranscodesWithURL(ctx context.Context, params CreateTranscodesWithURLParams) (CreateTranscodesResponse, int, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	req, err := http.NewRequest("POST", s.apiURL+"/wsgw/v1/mms/transcodes", bytes.NewBuffer(b))
	if err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	req.Header.Add("Authorization", "Basic "+s.auth)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	if res.StatusCode != 200 {
		output := CreateTranscodesErrorResponse{}
		if err := json.NewDecoder(res.Body).Decode(&output); err != nil {
			return CreateTranscodesResponse{}, res.StatusCode, err
		}

		return CreateTranscodesResponse{}, res.StatusCode, fmt.Errorf("mGage gave status code %d, response: %s", res.StatusCode, output.ErrorText)
	}

	output := CreateTranscodesResponse{}
	if err := json.NewDecoder(res.Body).Decode(&output); err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	return output, res.StatusCode, nil
}

func (s service) CreateTranscodesWithContent(ctx context.Context, params CreateTranscodesWithContentParams) (CreateTranscodesResponse, int, error) {
	payload := &bytes.Buffer{}
	w := multipart.NewWriter(payload)

	if err := w.WriteField("requestId", params.RequestID); err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	if err := w.WriteField("customerId", params.CustomerID); err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	if err := w.WriteField("mediaType", params.MediaType); err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	if err := w.WriteField("notificationUrl", params.NotificationURL); err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	if err := w.WriteField("Content-Type", params.ContentType); err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	file, err := w.CreateFormFile("filename", params.FileName)
	if err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	_, err = file.Write(params.Content)
	if err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	if err := w.Close(); err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	req, err := http.NewRequest("POST", s.apiURL+"/wsgw/v1/mms/transcodes", payload)
	if err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	req.Header.Add("Authorization", "Basic "+s.auth)
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	if res.StatusCode != 200 {
		output := CreateTranscodesErrorResponse{}
		if err := json.NewDecoder(res.Body).Decode(&output); err != nil {
			return CreateTranscodesResponse{}, res.StatusCode, err
		}

		return CreateTranscodesResponse{}, res.StatusCode, fmt.Errorf("mGage gave status code %d, response: %s", res.StatusCode, output.ErrorText)
	}

	output := CreateTranscodesResponse{}
	if err := json.NewDecoder(res.Body).Decode(&output); err != nil {
		return CreateTranscodesResponse{}, 0, err
	}

	return output, res.StatusCode, nil
}
