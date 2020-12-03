package rpc

import (
	"encoding/gob"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	"github.com/burstsms/mtmo-tp/backend/webhook/db"
)

const (
	Name = "Webhook"

	WebhookEventUnsubscribe = "unsubscribe"
	WebhookEventMMSStatus   = "mms_status"
	WebhookEventSMSStatus   = "sms_status"
	WebhookEventSMSInbound  = "sms_inbound"
	WebhookEventLinkHit     = "link_hit"
)

type Webhook struct {
	db   *db.DB
	name string
}

func NewWebhook(sdb *db.DB, name string) *Webhook {
	return &Webhook{
		db:   sdb,
		name: name,
	}
}

type Service struct {
	receiver *Webhook
}

func (s *Service) Name() string {
	return s.receiver.name
}

func (s *Service) Receiver() interface{} {
	return s.receiver
}

func NewService(sdb *db.DB) rpc.Service {
	gob.Register(time.Time{})
	gob.Register(map[string]interface{}{})
	return &Service{
		receiver: &Webhook{db: sdb, name: Name},
	}
}
