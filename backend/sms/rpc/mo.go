package rpc

import (
	"github.com/burstsms/mtmo-tp/backend/sms/worker/msg"
)

type QueueMOParams struct {
	MessageID     string
	Message       string
	To            string
	From          string
	SARID         string
	SARPartNumber string
	SARParts      string
}

func (s *SMSService) QueueMO(p QueueMOParams, r *NoReply) error {

	opts := RabbitPublishOptions{
		Exchange:     msg.MOMessage.Exchange,
		ExchangeType: msg.MOMessage.ExchangeType,
		RouteKey:     msg.MOMessage.RouteKey,
	}

	job := msg.MOMessageSpec{
		MessageID:     p.MessageID,
		Message:       p.Message,
		To:            p.To,
		From:          p.From,
		SARID:         p.SARID,
		SARPartNumber: p.SARPartNumber,
		SARParts:      p.SARParts,
	}

	err := s.db.Publish(job, opts)
	if err != nil {
		return err
	}

	return nil
}
