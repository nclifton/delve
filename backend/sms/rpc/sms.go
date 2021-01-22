package rpc

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	optOut "github.com/burstsms/mtmo-tp/backend/optout/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
	"github.com/burstsms/mtmo-tp/backend/sms/biz"
	"github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/sms/worker/msg"
	tracklink "github.com/burstsms/mtmo-tp/backend/track_link/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

func (s *SMSService) Send(p types.SendParams, r *types.SendReply) error {
	ctx := context.Background()
	uid := uuid.New().String()
	message := p.Message

	sender, err := s.senderRPC.FindByAddress(ctx, &senderpb.FindByAddressParams{
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
	err = biz.IsValidSender(sender.Sender, p.Sender, p.Country)
	if err != nil {
		return err
	}

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

	generateOptOutLinkReply, err := s.optOutRPC.GenerateOptOutLink(optOut.GenerateOptOutLinkParams{
		AccountID:   p.AccountID,
		MessageID:   uid,
		MessageType: "sms",
		Message:     message,
		Sender:      p.Sender,
	})
	if err != nil {
		return err
	}

	message = generateOptOutLinkReply.Message

	recipientNumber := p.Recipient
	var country string

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

	newSMS := types.SMS{
		ID:         uid,
		AccountID:  p.AccountID,
		MessageRef: p.MessageRef,
		Country:    country,
		Message:    message,
		SMSCount:   count,
		GSM:        isGSM,
		Recipient:  recipientNumber,
		Sender:     p.Sender,
		TrackLinks: p.TrackLinks,
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

func (s *SMSService) MarkSent(p types.MarkSentParams, r *types.NoReply) error {

	err := s.db.MarkSent(p.ID, p.MessageID)
	if err != nil {
		log.Printf("[Processing SMS] Could not mark sent SMS with ID: %s", p.ID)
		return err
	}

	// find sms by the given dlr messageid
	sms, err := s.db.FindSMSByID(p.ID, p.AccountID)
	if err != nil {
		log.Printf("[Processing SMS] SMS Not Found with ID: %s", p.ID)
		return err
	}

	// if it exists call the webhook service to send any status event webhooks
	_, err = s.webhookRPC.PublishSMSStatusUpdate(context.Background(), &webhookpb.PublishSMSStatusUpdateParams{
		AccountId:       sms.AccountID,
		SMSId:           sms.ID,
		MessageRef:      sms.MessageRef,
		Recipient:       sms.Recipient,
		Sender:          sms.Sender,
		Status:          `sent`,
		StatusUpdatedAt: timestamppb.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *SMSService) MarkFailed(p types.MarkFailedParams, r *types.NoReply) error {
	err := s.db.MarkFailed(p.ID)
	if err != nil {
		log.Printf("[Processing SMS] Could not mark failed SMS with ID: %s", p.ID)
		return err
	}

	// find sms by the given dlr messageid
	sms, err := s.db.FindSMSByID(p.ID, p.AccountID)
	if err != nil {
		log.Printf("[Processing SMS] SMS Not Found with ID: %s", p.ID)
		return err
	}

	// if it exists call the webhook service to send any status event webhooks
	_, err = s.webhookRPC.PublishSMSStatusUpdate(context.Background(), &webhookpb.PublishSMSStatusUpdateParams{
		AccountId:       sms.AccountID,
		SMSId:           sms.ID,
		MessageRef:      sms.MessageRef,
		Recipient:       sms.Recipient,
		Sender:          sms.Sender,
		Status:          `failed`,
		StatusUpdatedAt: timestamppb.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *SMSService) FindByID(p types.FindByIDParams, r *types.FindByIDReply) error {
	sms, err := s.db.FindSMSByID(p.ID, p.AccountID)
	if err != nil {
		return err
	}

	r.SMS = sms
	return nil
}
