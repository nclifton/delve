package processmo

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	rpc "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sms/worker/msg"
)

type SenderHandler struct {
	smsRPC *rpc.Client
}

func NewHandler(c *rpc.Client) *SenderHandler {
	return &SenderHandler{
		smsRPC: c,
	}
}

func (h *SenderHandler) OnFinalFailure(ctx context.Context, body []byte) error {
	return nil
}

func (h *SenderHandler) Handle(ctx context.Context, body []byte, headers map[string]interface{}) error {

	jobdata := &msg.MOMessageSpec{}
	err := json.NewDecoder(bytes.NewReader(body)).Decode(&jobdata)
	if err != nil {
		return rabbit.NewErrWorkerMessageParse(err.Error())
	}
	log.Printf("[MO process] got message: %+v", jobdata)

	// process the dlr
	err = h.smsRPC.ProcessMO(rpc.ProcessMOParams{
		MessageID:     jobdata.MessageID,
		Message:       jobdata.Message,
		To:            jobdata.To,
		From:          jobdata.From,
		SARID:         jobdata.SARID,
		SARPartNumber: jobdata.SARPartNumber,
		SARParts:      jobdata.SARParts,
	})
	if err != nil {
		log.Printf("[MO Process] error processing %s because: %s", jobdata.MessageID, err)
		return err
	}

	return nil
}
