package rpc

import (
	"time"

	"github.com/burstsms/mtmo-tp/backend/sms/worker/msg"
)

type QueueDLRParams struct {
	MessageID  string
	State      string
	ReasonCode string
	To         string
	Time       time.Time
	MCC        string
	MNC        string
}

func (s *SMSService) QueueDLR(p QueueDLRParams, r *NoReply) error {

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
