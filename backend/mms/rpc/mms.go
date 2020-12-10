package rpc

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/number"
	"github.com/burstsms/mtmo-tp/backend/mms/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/mms/worker"
	optOut "github.com/burstsms/mtmo-tp/backend/optout/rpc/client"
	tracklink "github.com/burstsms/mtmo-tp/backend/track_link/rpc/client"
	webhook "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
	"github.com/google/uuid"
)

func (s *MMSService) Send(p types.SendParams, r *types.SendReply) error {
	ctx := context.Background()

	uid := uuid.New().String()

	message := p.Message
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
		message = rsp.Message
	}

	generateOptOutLinkReply, err := s.svc.OptOut.GenerateOptOutLink(optOut.GenerateOptOutLinkParams{
		AccountID:   p.AccountID,
		MessageID:   uid,
		MessageType: "mms",
		Message:     message,
		Sender:      p.Sender,
	})
	if err != nil {
		return err
	}

	message = generateOptOutLinkReply.Message

	if len([]rune(p.Message)) > 1000 {
		return errors.New("message must be less than 1000 characters")
	}

	if len(p.ContentURLs) > 4 {
		return errors.New("you must provide no more then 4 content_urls")
	}

	recipientNumber := p.Recipient
	var country string

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

	newMMS := types.MMS{
		ID:          uid,
		AccountID:   p.AccountID,
		Subject:     p.Subject,
		Message:     message,
		Recipient:   recipientNumber,
		Sender:      p.Sender,
		Country:     country,
		MessageRef:  p.MessageRef,
		ContentURLs: p.ContentURLs,
		TrackLinks:  p.TrackLinks,
	}

	newMMS.Status = `pending`
	newMMS.ProviderKey = `fake`
	if p.ProviderKey != "" {
		newMMS.ProviderKey = p.ProviderKey
	}
	log.Printf("newMMS: %+v, Params: %+v", newMMS, p)

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

func (s *MMSService) UpdateStatus(p types.UpdateStatusParams, r *types.NoReply) error {
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

func (s *MMSService) FindByID(p types.FindByIDParams, r *types.FindByIDReply) error {
	ctx := context.Background()

	mms, err := s.db.FindByID(ctx, p.ID)
	if err != nil {
		return err
	}

	r.MMS = mms
	return nil
}
