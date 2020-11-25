package rpc

import (
	"context"
	"time"

	"github.com/burstsms/mtmo-tp/backend/mms/worker"
	webhook "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
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
	ShortenURLs bool      `json:"shorten_urls"`
	Unsub       bool      `json:"unsub"`
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
	ShortenURLs bool
}

type SendReply struct {
	MMS *MMS
}

func (s *MMSService) Send(p SendParams, r *SendReply) error {
	ctx := context.Background()

	newMMS := MMS{
		AccountID:   p.AccountID,
		Subject:     p.Subject,
		Message:     p.Message,
		Recipient:   p.Recipient,
		Sender:      p.Sender,
		Country:     p.Country,
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
