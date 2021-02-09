// +build integration

package client

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestCreateTranscodesWithURL(t *testing.T) {
	apiURL := os.Getenv("MGAGE_API_URL")
	user := os.Getenv("MGAGE_USER")
	password := os.Getenv("MGAGE_PASSWORD")
	service, err := NewService(apiURL, user, password)
	if err != nil {
		t.Fatal(err)
	}

	params := CreateTranscodesWithURLParams{
		CustomerID:      "7fb3a406-3005-47a9-a0ad-7a1840310ad2",
		RequestID:       "0b68158e-c284-46c4-bf51-a0b4f9e8a5f1",
		MediaType:       "IMAGE",
		MediaURL:        "http://www.acme.com/image.jpg",
		NotificationURL: "https://webhook.site/034c12e3-f3bc-4a08-afb1-1d221978cde9",
	}

	res, code, err := service.CreateTranscodesWithURL(context.Background(), params)
	if err != nil {
		if code != 0 {
			t.Fatalf("status code: %d, error: %s", code, err)
		}
		t.Fatal(err)
	}

	log.Printf("RESULT: %+v, status code: %d", res, code)
}

func TestCreateTranscodesWithContent(t *testing.T) {
	apiURL := os.Getenv("MGAGE_API_URL")
	user := os.Getenv("MGAGE_USER")
	password := os.Getenv("MGAGE_PASSWORD")
	service, err := NewService(apiURL, user, password)
	if err != nil {
		t.Fatal(err)
	}

	content, err := ioutil.ReadFile("./test/image-test.png")
	if err != nil {
		t.Fatal(err)
	}

	params := CreateTranscodesWithContentParams{
		CustomerID:      "7fb3a406-3005-47a9-a0ad-7a1840310ad2",
		RequestID:       "0b68158e-c284-46c4-bf51-a0b4f9e8a5e5",
		MediaType:       "IMAGE",
		NotificationURL: "https://webhook.site/034c12e3-f3bc-4a08-afb1-1d221978cde9",
		ContentType:     "image/png",
		FileName:        "image-test.png",
		Content:         content,
	}

	res, code, err := service.CreateTranscodesWithContent(context.Background(), params)
	if err != nil {
		if code != 0 {
			t.Fatalf("status code: %d, error: %s", code, err)
		}
		t.Fatal(err)
	}

	log.Printf("RESULT: %+v, status code: %d", res, code)
}
