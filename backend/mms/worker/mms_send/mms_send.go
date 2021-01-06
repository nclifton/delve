package mms_send

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"
	mmsRPC "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/mms/worker"
)

// TODO: to fix ? shouldn't be exported or in a common folder?
const (
	MMSStatusNew       = "new"
	MMSStatusProcessed = "processed"
	MMSStatusFailed    = "failed"
	MMSStatusSent      = "sent"
)

type mm7RPCClient interface {
	Send(p mm7RPC.SendParams) error
}

type mmsRPCClient interface {
	UpdateStatus(p mmsRPC.UpdateStatusParams) error
}

type MMSSendHandler struct {
	mm7RPC mm7RPCClient
	mmsRPC mmsRPCClient
	log    *logger.StandardLogger
}

func NewHandler(mm7c mm7RPCClient, mmsc mmsRPCClient) *MMSSendHandler {
	return &MMSSendHandler{
		mm7RPC: mm7c,
		mmsRPC: mmsc,
		log:    logger.NewLogger(),
	}
}

func (h *MMSSendHandler) OnFinalFailure(ctx context.Context, body []byte) error {
	return nil
}

func (h *MMSSendHandler) Handle(ctx context.Context, body []byte, headers map[string]interface{}) error {
	msg := &worker.Job{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&msg); err != nil {
		h.logError(ctx, msg, "", err.Error(), "Decoding job failed")
		return err
	}

	// send to mm7 service
	if err := h.mm7RPC.Send(mm7RPC.SendParams{
		ID:          msg.ID,
		Subject:     msg.Subject,
		Message:     msg.Message,
		Sender:      msg.Sender,
		Recipient:   msg.Recipient,
		ContentURLs: msg.ContentURLs,
		ProviderKey: msg.ProviderKey,
	}); err != nil {
		h.logError(ctx, msg, MMSStatusFailed, err.Error(), "Problem sending to mm7")

		return h.mmsRPC.UpdateStatus(mmsRPC.UpdateStatusParams{
			ID:          msg.ID,
			Status:      MMSStatusFailed,
			Description: err.Error(),
		})
	}

	h.logSuccess(ctx, msg, MMSStatusSent, "", "MMS send successfully to mm7 service")

	return nil
}

func (h *MMSSendHandler) logError(ctx context.Context, msg *worker.Job, status, description, label string) {
	fields := logger.Fields{
		"ID":          msg.ID,
		"Sender":      msg.Sender,
		"Recipient":   msg.Recipient,
		"Status":      status,
		"Description": description,
	}
	h.log.Fields(ctx, fields).Error(label)
}

func (h *MMSSendHandler) logSuccess(ctx context.Context, msg *worker.Job, status, description, label string) {
	fields := logger.Fields{
		"ID":          msg.ID,
		"Sender":      msg.Sender,
		"Recipient":   msg.Recipient,
		"Status":      status,
		"Description": description,
	}

	h.log.Fields(ctx, fields).Info(label)
}
