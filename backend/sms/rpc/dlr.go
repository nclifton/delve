package rpc

import (
	"log"
	"time"

	"github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/sms/worker/msg"
	webhookRPC "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
)

func (s *SMSService) QueueDLR(p types.QueueDLRParams, r *types.NoReply) error {

	opts := RabbitPublishOptions{
		Exchange:     msg.DLRMessage.Exchange,
		ExchangeType: msg.DLRMessage.ExchangeType,
		RouteKey:     msg.DLRMessage.RouteKey,
	}

	job := msg.DLRMessageSpec{
		MessageID:  p.MessageID,
		State:      p.State,
		ReasonCode: p.ReasonCode,
		To:         p.To,
		Time:       p.Time,
		MCC:        p.MCC,
		MNC:        p.MNC,
	}

	err := s.db.Publish(job, opts)
	if err != nil {
		return err
	}

	return nil
}

func (s *SMSService) ProcessDLR(p types.ProcessDLRParams, r *types.NoReply) error {
	// find sms by the given dlr messageid
	sms, err := s.db.FindSMSByMessageID(p.MessageID)
	if err != nil {
		log.Printf("[Processing DLR] SMS Not Found with MessageID: %s", p.MessageID)
		return err
	}

	log.Printf("[Processing DLR] Found SMS: %+v", sms)

	// update the sms status with the dlr status
	err = s.db.MarkStatus(sms.ID, p.State)
	if err != nil {
		return err
	}

	// if it exists call the webhook service to send any status event webhooks
	err = s.webhookRPC.PublishSMSStatusUpdate(webhookRPC.PublishSMSStatusUpdateParams{
		AccountID:       sms.AccountID,
		SMSID:           sms.ID,
		MessageRef:      sms.MessageRef,
		Recipient:       sms.Recipient,
		Sender:          sms.Sender,
		Status:          p.State,
		StatusUpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}
