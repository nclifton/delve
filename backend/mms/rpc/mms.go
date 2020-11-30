package rpc

import (
	"context"
	"errors"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/number"
	"github.com/burstsms/mtmo-tp/backend/mms/worker"
	tracklink "github.com/burstsms/mtmo-tp/backend/track_link/rpc/client"
	webhook "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
	"github.com/google/uuid"
)

type MMS struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ProviderKey string    `json:"provider_key"`
	MessageID   string    `json:"message_id"`
	MessageRef  string    `json:"message_ref"`
	Country     string    `json:"country"`
	Subject     string    `json:"subject"`
	Message     string    `json:"message"`
	ContentURLs []string  `json:"content_urls"`
	Recipient   string    `json:"recipient"`
	Sender      string    `json:"sender"`
	Status      string    `json:"status"`
	TrackLinks  bool      `json:"track_links"`
}

type SendParams struct {
	AccountID   string
	Subject     string
	Message     string
	Recipient   string
	Sender      string
	Country     string
	MessageRef  string
	ContentURLs []string
	TrackLinks  bool
}

type SendReply struct {
	MMS *MMS
}

func (s *MMSService) Send(p SendParams, r *SendReply) error {
	ctx := context.Background()

	uid := uuid.New().String()

	msg := p.Message
	if p.TrackLinks {
		rsp, err := s.svc.TrackLink.GenerateTrackLinks(tracklink.GenerateTrackLinksParams{
			AccountID:   p.AccountID,
			MessageID:   uid,
			MessageType: Name,
			Message:     p.Message,
		})
		if err != nil {
			return err
		}
		msg = rsp.Message
	}

	if len([]rune(p.Message)) > 1000 {
		return errors.New("message must be less than 1000 characters")
	}

	if len(p.ContentURLs) > 4 {
		return errors.New("you must provide no more then 4 content_urls")
	}

	recipientNumber := p.Recipient
	var country string
	var err error

	if p.Country != "" {
		recipientNumber, country, err = number.ParseMobileCountry(recipientNumber, p.Country)
		if err != nil {
			return err
		}
	} else {
		country, err = number.GetCountryFromPhone(recipientNumber)
		if err != nil {
			return err
		}
	}

	newMMS := MMS{
		ID:          uid,
		AccountID:   p.AccountID,
		Subject:     p.Subject,
		Message:     msg,
		Recipient:   recipientNumber,
		Sender:      p.Sender,
		Country:     country,
		MessageRef:  p.MessageRef,
		ContentURLs: p.ContentURLs,
	}

	newMMS.Status = `pending`
	newMMS.ProviderKey = `fake`

	mms, err := s.db.InsertMMS(ctx, newMMS)
	if err != nil {
		return err
	}
	r.MMS = mms

	job := worker.Job{
		ID:          mms.ID,
		AccountID:   mms.AccountID,
		Sender:      mms.Sender,
		Subject:     mms.Subject,
		ContentURLs: mms.ContentURLs,
		Recipient:   mms.Recipient,
		ProviderKey: mms.ProviderKey,
		Message:     mms.Message,
	}

	err = s.db.Publish(job, worker.MMSSendQueueName)

	return err
}

type FindByIDParams struct {
	ID        string
	AccountID string
}

type FindByIDReply struct {
	MMS *MMS
}

func (s *MMSService) FindByID(p FindByIDParams, r *FindByIDReply) error {
	ctx := context.Background()

	mms, err := s.db.FindByIDAndAccountID(ctx, p.ID, p.AccountID)
	if err != nil {
		return err
	}

	r.MMS = mms
	return nil
}

type UpdateStatusParams struct {
	ID          string
	MessageID   string
	Status      string
	Description string
}

func (s *MMSService) UpdateStatus(p UpdateStatusParams, r *NoReply) error {
	ctx := context.Background()

	mms, err := s.db.FindByID(ctx, p.ID)
	if err != nil {
		return err
	}

	if err := s.db.UpdateStatus(ctx, p.ID, p.MessageID, p.Status); err != nil {
		return err
	}

	return s.svc.Webhook.PublishMMSStatusUpdate(webhook.PublishMMSStatusUpdateParams{
		AccountID:         mms.AccountID,
		MMSID:             mms.ID,
		MessageRef:        mms.MessageRef,
		Recipient:         mms.Recipient,
		Sender:            mms.Sender,
		Status:            p.Status,
		StatusDescription: p.Description,
		StatusUpdatedAt:   time.Now(),
	})
}
