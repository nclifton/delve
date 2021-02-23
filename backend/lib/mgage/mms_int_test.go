// +build integration

package mgage

import (
	"context"
	"log"
	"os"
	"testing"
)

func TestSendMMS(t *testing.T) {
	apiURL := os.Getenv("MGAGE_API_URL")
	user := os.Getenv("MGAGE_USER")
	password := os.Getenv("MGAGE_PASSWORD")
	service, err := NewService(apiURL, user, password)
	if err != nil {
		t.Fatal(err)
	}

	params := SendMMSParams{
		CustomerID:      "7fb3a406-3005-47a9-a0ad-7a1840310ad2",
		RequestID:       "0b68158e-c284-46c4-bf51-a0b4f9e8a5f1",
		From:            "123-667-000",
		To:              []string{"5157791970"},
		Subject:         "mt subject",
		RefMessageID:    "mo_message_id",
		ReportingKey1:   "key_1",
		ReportingKey2:   "key_2",
		AllowAdaptation: "true",
		ForwardLock:     "true",
		Messages: []Message{
			{
				Text:        "Hello mGage!",
				DisplayName: "img1.jpg",
				MediaType:   "IMAGE",
				MediaURL:    "http://www.acme.com/image.jpg",
			},
		},
	}

	res, code, err := service.SendMMS(context.Background(), params)
	if err != nil {
		if code != 0 {
			t.Fatalf("status code: %d, error: %s", code, err)
		}
		t.Fatal(err)
	}

	log.Printf("RESULT: %+v, status code: %d", res, code)
}
