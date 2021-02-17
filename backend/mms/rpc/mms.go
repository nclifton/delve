package rpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	"github.com/burstsms/mtmo-tp/backend/lib/number"
	"github.com/burstsms/mtmo-tp/backend/mms/biz"
	"github.com/burstsms/mtmo-tp/backend/mms/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/mms/worker"
	optOut "github.com/burstsms/mtmo-tp/backend/optout/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
	tracklink "github.com/burstsms/mtmo-tp/backend/track_link/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

func (s *MMSService) Send(p types.SendParams, r *types.SendReply) error {
	ctx := context.Background()
	uid := uuid.New().String()

	sender, err := s.svc.Sender.FindSenderByAddressAndAccountID(ctx, &senderpb.FindSenderByAddressAndAccountIDParams{
		AccountId: p.AccountID,
		Address:   p.Sender,
	})
	if err != nil {
		errStatus, ok := status.FromError(err)
		if ok && errStatus.Code() == codes.NotFound {
			return errorlib.ErrInvalidSenderNotFound
		}
		return err
	}

	if err := biz.IsValidSender(sender.Sender, p.Sender, p.Country); err != nil {
		return err
	}

	if len([]rune(p.Message)) > 1000 {
		return errorlib.ErrInvalidMMSLengthMessage
	}

	if len(p.ContentURLs) > 4 {
		return errorlib.ErrInvalidMMSLengthContentURLs
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
			return errorlib.ErrInvalidRecipientInternationalNumber

		}
	}

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
		Status:      "pending",
		ProviderKey: sender.Sender.GetMMSProviderKey(),
	}

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

	_, err = s.svc.Webhook.PublishMMSStatusUpdate(ctx, &webhookpb.PublishMMSStatusUpdateParams{
		AccountId:         mms.AccountID,
		MMSId:             mms.ID,
		MessageRef:        mms.MessageRef,
		Recipient:         mms.Recipient,
		Sender:            mms.Sender,
		Status:            p.Status,
		StatusDescription: p.Description,
		StatusUpdatedAt:   timestamppb.New(time.Now()),
	})

	return err
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
