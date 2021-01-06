package post

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/msg"
)

type Webhook struct {
	log     *logger.StandardLogger
	client  *http.Client
	limiter Limiter
}

type Limiter interface {
	Allow(url string, rate float64, burst int) bool
}

type Publisher interface {
	Publish(msg interface{}) error
}

func NewHandler(client *http.Client, limiter Limiter) *Webhook {
	return &Webhook{logger.NewLogger(), client, limiter}
}

func (h *Webhook) OnFinalFailure(ctx context.Context, body []byte) error {
	return nil
}

func (h *Webhook) Handle(ctx context.Context, body []byte, headers map[string]interface{}) error {
	data := &msg.WebhookMessageSpec{}

	err := json.NewDecoder(bytes.NewReader(body)).Decode(&data)
	if err != nil {
		h.log.Error(ctx, "json.NewDecoder", err.Error())
		return rabbit.NewErrWorkerMessageParse(err.Error())
	}

	if data.RateLimit > 0 && !h.limiter.Allow(data.URL, float64(data.RateLimit), data.RateLimit) {
		return errors.New("Hit Ratelimit")
	}

	payload, err := json.Marshal(data.Payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", data.URL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return rabbit.NewErrRetryWorkerMessage(fmt.Sprintf("Failed sending webhook to: %s With params: %+v Error: %s", data.URL, data.Payload, err.Error()))
	}
	defer resp.Body.Close()

	return nil
}
