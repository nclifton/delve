package rpc

import (
	"time"

	"github.com/burstsms/mtmo-tp/backend/sms/biz"
	"github.com/burstsms/mtmo-tp/backend/sms/worker/msg"
	tracklink "github.com/burstsms/mtmo-tp/backend/track_link/rpc/client"
	"github.com/google/uuid"
)

type SMS struct {
	ID         string
	AccountID  string
	MessageID  string
	Recipient  string
	Sender     string
	Country    string
	MessageRef string
	Message    string
	Status     string
	SMSCount   int
	GSM        bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	TrackLinks bool
}

type SendParams struct {
	AccountID  string
	Message    string
	Recipient  string
	Sender     string
	Country    string
	MessageRef string
	AlarisUser string
	AlarisPass string
	AlarisURL  string
	TrackLinks bool
}

type SendReply struct {
	SMS *SMS
}

func (s *SMSService) Send(p SendParams, r *SendReply) error {
	uid := uuid.New().String()

	message := p.Message
	if p.TrackLinks {
		rsp, err := s.tracklinkRPC.GenerateTrackLinks(tracklink.GenerateTrackLinksParams{
			AccountID:   p.AccountID,
			MessageID:   uid,
			MessageType: Name,
			Message:     p.Message,
		})
		if err != nil {
			return err
		}
		message = rsp.Message
	}

	recipientNumber := p.Recipient
	var country string
	var err error

	if p.Country != "" {
		recipientNumber, country, err = biz.ParseMobileCountry(recipientNumber, p.Country)
		if err != nil {
			return err
		}
	} else {
		country, err = biz.GetCountryFromPhone(recipientNumber)
		if err != nil {
			return err
		}
	}

	options := biz.SMSOptions{
		MaxParts:         4,
		OptOutLinkDomain: s.features.OptOutLinkDomain,
		TrackLinkDomain:  s.features.TrackLinkDomain,
		TrackLink:        p.TrackLinks,
	}
	// check the sms size
	count, err := biz.IsValidSMS(p.Message, options)
	if err != nil {
		return err
	}

	// check if its a GSM compat message
	isGSM := biz.IsGSMString(p.Message)

	newSMS := SMS{
		ID:         uid,
		AccountID:  p.AccountID,
		MessageRef: p.MessageRef,
		Country:    country,
		Message:    message,
		SMSCount:   count,
		GSM:        isGSM,
		Recipient:  recipientNumber,
		Sender:     p.Sender,
	}

	sms, err := s.db.InsertSMS(newSMS)
	if err != nil {
		return err
	}

	opts := RabbitPublishOptions{
		Exchange:     msg.SMSSendMessage.Exchange,
		ExchangeType: msg.SMSSendMessage.ExchangeType,
		RouteKey:     msg.SMSSendMessage.RouteKey,
	}

	job := msg.SMSSendMessageSpec{
		ID:         sms.ID,
		Recipient:  sms.Recipient,
		Sender:     sms.Sender,
		Message:    sms.Message,
		AccountID:  sms.AccountID,
		AlarisUser: p.AlarisUser,
		AlarisPass: p.AlarisPass,
		AlarisURL:  p.AlarisURL,
	}

	err = s.db.Publish(job, opts)
	if err != nil {
		return err
	}

	r.SMS = sms
	return nil
}

type MarkSentParams struct {
	ID        string
	MessageID string
}

func (s *SMSService) MarkSent(p MarkSentParams, r *NoReply) error {
	return s.db.MarkSent(p.ID, p.MessageID)
}

type MarkFailedParams struct {
	ID string
}

func (s *SMSService) MarkFailed(p MarkFailedParams, r *NoReply) error {
	return s.db.MarkFailed(p.ID)
}
