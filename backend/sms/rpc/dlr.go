package rpc

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/sms/worker/msg"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

const (
	DLRStatusDelivered   = "delivered"
	DLRStatusSoftBounce  = "soft_bounce"
	DLRStatusHardBounce  = "hard_bounce"
	DLRStatusAccepted    = "accepted"
	DLRStatusNotAccepted = "not_accepted"
	DLRStatusPending     = "pending"

	dlrCodeDelivered   = "DELIVRD"
	dlrCodeAccepted    = "ACCEPTD"
	dlrCodeExpired     = "EXPIRED"
	dlrCodeDeleted     = "DELETED"
	dlrCodeUndelivered = "UNDELIV"
	dlrCodeRejected    = "REJECTD"
	//dlrCodeSystemExpired = "EXPIRED"
	dlrCodeEnroute = "ENROUTE"
	dlrCodeSent    = "SENT"
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

	var status string

	switch p.State {
	case dlrCodeDelivered:
		status = DLRStatusDelivered
	case dlrCodeExpired:
		fallthrough
	case dlrCodeDeleted:
		fallthrough
	case dlrCodeRejected:
		fallthrough
	case dlrCodeUndelivered:
		status = DLRStatusHardBounce
	case dlrCodeEnroute:
		fallthrough
	case dlrCodeAccepted:
		fallthrough
	case dlrCodeSent:
		status = DLRStatusAccepted
	default:
		status = p.State
	}

	// update the sms status with the dlr status
	err = s.db.MarkStatus(sms.ID, status)
	if err != nil {
		return err
	}

	// if it exists call the webhook service to send any status event webhooks
	_, err = s.webhookRPC.PublishSMSStatusUpdate(context.Background(), &webhookpb.PublishSMSStatusUpdateParams{
		AccountId:       sms.AccountID,
		SMSId:           sms.ID,
		MessageRef:      sms.MessageRef,
		Recipient:       sms.Recipient,
		Sender:          sms.Sender,
		Status:          status,
		StatusUpdatedAt: timestamppb.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}
