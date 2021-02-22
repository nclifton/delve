package optussubmitworker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	tcl "github.com/burstsms/mtmo-tp/backend/lib/optus/client"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/mm7/worker"
)

type mm7RPCClient interface {
	UpdateStatus(p mm7RPC.UpdateStatusParams) error
	CheckRateLimit(p mm7RPC.CheckRateLimitParams) (r *mm7RPC.CheckRateLimitReply, err error)
	GetCachedContent(p mm7RPC.GetCachedContentParams) (r *mm7RPC.GetCachedContentReply, err error)
}

type optusClient interface {
	PostMM7(params tcl.PostMM7Params, soaptmpl *template.Template) (tcl.PostMM7Response, int, error)
}

type OptusSubmitHandler struct {
	mm7RPC   mm7RPCClient
	optus    optusClient
	log      *logger.StandardLogger
	soaptmpl *template.Template
}

func NewHandler(c mm7RPCClient, optus optusClient, soaptmpl *template.Template) *OptusSubmitHandler {
	return &OptusSubmitHandler{
		mm7RPC:   c,
		optus:    optus,
		log:      logger.NewLogger(),
		soaptmpl: soaptmpl,
	}
}

func (h *OptusSubmitHandler) OnFinalFailure(ctx context.Context, body []byte) error {
	return nil
}

const (
	MMSStatusFailed = "failed"
	MMSStatusSent   = "sent"
)

func (h *OptusSubmitHandler) Handle(ctx context.Context, body []byte, headers map[string]interface{}) error {
	msg := &worker.SubmitMessage{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&msg); err != nil {
		h.logError(ctx, msg, "", err.Error(), "Decoding job failed")
		return err
	}

	r, err := h.mm7RPC.CheckRateLimit(mm7RPC.CheckRateLimitParams{ProviderKey: worker.FakeProviderKey})
	if err != nil {
		h.logError(ctx, msg, "", err.Error(), "Unexpected mm7RPC.CheckRateLimit response")
		return err
	}

	if !r.Allow {
		return rabbit.NewErrRetryWorkerMessage(fmt.Sprintf("Failed sending message id: %s Error: rate limit reached", msg.ID))
	}

	var images [][]byte
	for _, url := range msg.ContentURLs {
		r, err := h.mm7RPC.GetCachedContent(mm7RPC.GetCachedContentParams{
			ContentURL: url,
		})
		if err != nil {
			h.logError(ctx, msg, "", err.Error(), "Unexpected mm7RPC.GetCachedContent response")
			return err
		}

		images = append(images, r.Content)
	}

	var status string
	var description string

	result, _, err := h.optus.PostMM7(tcl.PostMM7Params{
		ID:        msg.ID,
		Subject:   msg.Subject,
		Message:   msg.Message,
		Sender:    msg.Sender,
		Recipient: msg.Recipient,
		Images:    images,
	}, h.soaptmpl)
	if err != nil {
		status = MMSStatusFailed
		description = err.Error()
		h.logError(ctx, msg, status, description, "Optus http request failed")
		err := h.updateStatus(msg.ID, status, description)
		return err
	}

	description = result.Body.SubmitRsp.Status.StatusText

	if result.Body.SubmitRsp.Status.StatusCode != "1000" {
		status = MMSStatusFailed
		h.logError(ctx, msg, status, description, "Received error status from Tecloo")
		err := h.updateStatus(msg.ID, status, description)
		return err
	}

	status = MMSStatusSent
	h.logSuccess(ctx, msg, status, description, "Optus Submit Worker Successful send")
	err = h.updateStatus(msg.ID, status, description)
	return err
}

func (h *OptusSubmitHandler) updateStatus(id, status, description string) error {
	err := h.mm7RPC.UpdateStatus(mm7RPC.UpdateStatusParams{
		ID:          id,
		Status:      status,
		Description: description,
	})

	return err
}

func (h *OptusSubmitHandler) logError(ctx context.Context, msg *worker.SubmitMessage, status, description, label string) {
	fields := logger.Fields{
		"ID":          msg.ID,
		"Sender":      msg.Sender,
		"Recipient":   msg.Recipient,
		"Status":      status,
		"Description": description,
	}

	h.log.Fields(ctx, fields).Error(label)
}

func (h *OptusSubmitHandler) logSuccess(ctx context.Context, msg *worker.SubmitMessage, status, description, label string) {
	fields := logger.Fields{
		"ID":          msg.ID,
		"Sender":      msg.Sender,
		"Recipient":   msg.Recipient,
		"Status":      status,
		"Description": description,
	}

	h.log.Fields(ctx, fields).Info(label)
}
