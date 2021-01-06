package fakemm7submitworker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	tcl "github.com/burstsms/mtmo-tp/backend/lib/tecloo/client"
	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/mm7/worker"
)

type mm7RPCClient interface {
	UpdateStatus(p mm7RPC.UpdateStatusParams) error
	ProviderSpec(p mm7RPC.ProviderSpecParams) (r *mm7RPC.ProviderSpecReply, err error)
	CheckRateLimit(p mm7RPC.CheckRateLimitParams) (r *mm7RPC.CheckRateLimitReply, err error)
	GetCachedContent(p mm7RPC.GetCachedContentParams) (r *mm7RPC.GetCachedContentReply, err error)
}

type teclooSvc interface {
	PostMM7(params tcl.PostMM7Params, soaptmpl *template.Template) (tcl.PostMM7Response, int, error)
}

type FakeMM7SubmitHandler struct {
	mm7RPC   mm7RPCClient
	tecloo   teclooSvc
	log      *logger.StandardLogger
	soaptmpl *template.Template
}

func NewHandler(c mm7RPCClient, tecloo teclooSvc, soaptmpl *template.Template) *FakeMM7SubmitHandler {
	return &FakeMM7SubmitHandler{
		mm7RPC:   c,
		tecloo:   tecloo,
		log:      logger.NewLogger(),
		soaptmpl: soaptmpl,
	}
}

func (h *FakeMM7SubmitHandler) OnFinalFailure(context.Context, []byte) error {
	return nil
}

const (
	MMSStatusFailed = "failed"
	MMSStatusSent   = "sent"
)

func (h *FakeMM7SubmitHandler) Handle(ctx context.Context, body []byte, headers map[string]interface{}) error {
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

	psReply, err := h.mm7RPC.ProviderSpec(mm7RPC.ProviderSpecParams{
		ProviderKey: worker.FakeProviderKey,
	})
	if err != nil {
		h.logError(ctx, msg, "", err.Error(), "Unexpected mm7RPC.ProviderSpec response")
		return err
	}

	var images [][]byte
	for _, url := range msg.ContentURLs {
		r, err := h.mm7RPC.GetCachedContent(mm7RPC.GetCachedContentParams{
			ContentURL: url,
		})
		if err != nil {
			h.logError(ctx, msg, "", err.Error(), "Unexpected mm7RPC.GetCachedContent response")
			return h.updateStatus(msg.ID, "", MMSStatusFailed, err.Error())
		}

		// validation image size
		if len(r.Content) > psReply.ImageSizeMaxKB*1000 {
			description := fmt.Sprintf("Total image size > %dkb", psReply.ImageSizeMaxKB)
			h.logError(ctx, msg, MMSStatusFailed, description, "Validation error")
			return h.updateStatus(msg.ID, "", MMSStatusFailed, description)
		}

		images = append(images, r.Content)
	}

	var status string
	var description string

	result, _, err := h.tecloo.PostMM7(tcl.PostMM7Params{
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
		h.logError(ctx, msg, status, description, "Tecloo http request failed")
		err := h.updateStatus(msg.ID, "", status, description)
		return err
	}

	description = result.Body.SubmitRsp.Status.StatusText

	if result.Body.SubmitRsp.Status.StatusCode != "1000" {
		status = MMSStatusFailed
		h.logError(ctx, msg, status, description, "Received error status from Tecloo")
		err := h.updateStatus(msg.ID, result.Body.SubmitRsp.MessageID, status, description)
		return err
	}

	status = MMSStatusSent
	h.logSuccess(ctx, msg, status, description, "Fake MM7 Submit Worker Successful send")

	return h.updateStatus(msg.ID, result.Body.SubmitRsp.MessageID, status, description)
}

func (h *FakeMM7SubmitHandler) updateStatus(id, messageID, status, description string) error {
	return h.mm7RPC.UpdateStatus(mm7RPC.UpdateStatusParams{
		ID:          id,
		MessageID:   messageID,
		Status:      status,
		Description: description,
	})
}

func (h *FakeMM7SubmitHandler) logError(ctx context.Context, msg *worker.SubmitMessage, status, description, label string) {
	fields := logger.Fields{
		"ID":          msg.ID,
		"Sender":      msg.Sender,
		"Recipient":   msg.Recipient,
		"Status":      status,
		"Description": description,
	}

	h.log.Fields(ctx, fields).Error(label)
}

func (h *FakeMM7SubmitHandler) logSuccess(ctx context.Context, msg *worker.SubmitMessage, status, description, label string) {
	fields := logger.Fields{
		"ID":          msg.ID,
		"Sender":      msg.Sender,
		"Recipient":   msg.Recipient,
		"Status":      status,
		"Description": description,
	}

	h.log.Fields(ctx, fields).Info(label)
}
