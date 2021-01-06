package fakemm7dldrworker

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"
	"text/template"

	"github.com/burstsms/mtmo-tp/backend/lib/mm7utils"
)

func TestHandle(t *testing.T) {
	t.Run("returns an error if job is invalid", func(t *testing.T) {
		client := MockRPCClient{}
		w := NewHandler(client)
		body := []byte("Some payload")

		if err := w.Handle(context.Background(), body, nil); err == nil {
			t.Error(err)
		}
	})

	t.Run("returns an error if request is not dl or dr", func(t *testing.T) {
		client := MockRPCClient{}
		w := NewHandler(client)

		body := []byte("Some request body")
		headers := map[string]interface{}{
			"Content-Type": "application/json",
		}

		if err := w.Handle(context.Background(), body, headers); err == nil {
			t.Error(err)
		}
	})

	t.Run("calls mm7 for DR requests", func(t *testing.T) {
		client := MockRPCClient{}
		w := NewHandler(client)

		soaptmpl, err := template.New("tecloo_submit_dr.soap.tmpl").ParseFiles("./test_files/tecloo_submit_dr.soap.tmpl")
		if err != nil {
			t.Error(err)
		}

		params := mm7utils.DeliveryReportParams{
			MessageID:     "3000",
			StatusCode:    "1000",
			StatusText:    "statustext",
			Recipient:     "61422265707",
			Sender:        "61422265777",
			TransactionID: "2000",
		}

		soapdata, contentType, err := mm7utils.GenerateMM7DeliveryReport(params, soaptmpl)
		if err != nil {
			t.Error(err)
		}

		body, err := ioutil.ReadAll(soapdata)
		if err != nil {
			t.Error(err)
		}

		headers := map[string]interface{}{
			"Content-Type": contentType,
		}

		if err = w.Handle(context.Background(), body, headers); err != nil {
			t.Error(err)
		}
	})

	t.Run("returns an error if mms image fails to be cached", func(t *testing.T) {
		client := MockRPCClient{
			Error: errors.New("Failed to store image"),
		}

		w := NewHandler(client)

		soaptmpl, err := template.New("test").Parse(`{{.TransactionID}}`)
		if err != nil {
			t.Error(err)
		}

		params := mm7utils.DeliverRequestParams{
			TransactionID:    "1000",
			Subject:          "My MMS test message",
			Priority:         "Normal",
			Sender:           "61455123456",
			Recipient:        fmt.Sprintf("6142226%s", "1000"),
			AllowAdaptations: true,
		}

		text := "Test content"

		soapdata, contentType, err := mm7utils.GenerateMM7DeliverRequest(params, soaptmpl, text, [][]byte{[]byte("Test image content 1")})
		if err != nil {
			t.Error(err)
		}

		body, err := ioutil.ReadAll(soapdata)
		if err != nil {
			t.Error(err)
		}

		headers := map[string]interface{}{
			"Content-Type": contentType,
		}

		if err = w.Handle(context.Background(), body, headers); err == nil {
			t.Error(err)
		}
	})

	t.Run("calls mm7 for DL requests", func(t *testing.T) {
		client := MockRPCClient{}
		w := NewHandler(client)

		soaptmpl, err := template.New("tecloo_delivery.soap.tmpl").ParseFiles("./test_files/tecloo_delivery.soap.tmpl")
		if err != nil {
			t.Error(err)
		}

		params := mm7utils.DeliverRequestParams{
			TransactionID:    "123",
			Subject:          "My MMS test message",
			Priority:         "Normal",
			Sender:           "61455123456",
			Recipient:        fmt.Sprintf("6142226%s", "1000"),
			AllowAdaptations: true,
		}

		image1, err := ioutil.ReadFile("./test_files/image1.png")
		if err != nil {
			t.Error(err)
		}

		text := "Test content"

		soapdata, contentType, err := mm7utils.GenerateMM7DeliverRequest(params, soaptmpl, text, [][]byte{image1})
		if err != nil {
			t.Error(err)
		}

		body, err := ioutil.ReadAll(soapdata)
		if err != nil {
			t.Error(err)
		}

		headers := map[string]interface{}{
			"Content-Type": contentType,
		}

		if err = w.Handle(context.Background(), body, headers); err != nil {
			t.Error(err)
		}
	})
}
