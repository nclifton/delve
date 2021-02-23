package mgage

import "time"

type TranscodesPOSTRequest struct {
	RequestID   string `json:"requestId"`
	Type        string `json:"type"`
	TranscodeID string `json:"transcodeId"`
	Status      string `json:"status"`
}

type TranscodesPOSTResponse struct {
	Status string `json:"status"`
}

type MTPOSTRequest struct {
	RequestID     string    `json:"requestId"`
	Type          string    `json:"type"`
	From          string    `json:"from"`
	To            string    `json:"to"`
	MessageID     string    `json:"messageId"`
	SentAt        time.Time `json:"sentAt"`
	ReportingKey1 string    `json:"reportingKey1"`
	ReportingKey2 string    `json:"reportingKey2"`
	ErrorText     string    `json:"errorText"`
}

type MTPOSTResponse struct {
	Status string `json:"status"`
}

type MOPOSTRequest struct {
	Type      string    `json:"type"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	MessageID string    `json:"messageId"`
	SentAt    time.Time `json:"sentAt"`
	Message   MOMessage `json:"message"`
}

type MOMessage struct {
	Slide     string `json:"slide"`
	Text      string `json:"text"`
	MediaType string `json:"mediaType"`
	MediaURL  string `json:"mediaUrl"`
}

type MOPOSTResponse struct {
	Status string `json:"status"`
}
