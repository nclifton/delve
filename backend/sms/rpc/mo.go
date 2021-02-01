package rpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	optOutRPC "github.com/burstsms/mtmo-tp/backend/optout/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
	"github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/sms/worker/msg"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

func (s *SMSService) QueueMO(p types.QueueMOParams, r *types.NoReply) error {

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

var ErrInsufficientParts = errors.New("Insufficient parts to combine MO")

func (s *SMSService) checkMultiPart(p *types.ProcessMOParams) error {
	// Check for multipart
	if !strings.Contains(p.SARID, `sarId`) {

		// Store the part
		err := s.db.StoreSMSPart(p.SARID, p.MessageID, p.Message, p.SARPartNumber)
		if err != nil {
			return err
		}

		// count the parts
		count, err := s.db.CountStoredParts(p.SARID)
		if err != nil {
			return err
		}
		parcount, err := strconv.ParseInt(p.SARParts, 10, 64)
		if err != nil {
			return err
		}

		// do we have all the parts
		if count == parcount {
			// Join the message into one record
			parts, err := s.db.GetAllSMSParts(p.SARID)
			if err != nil {
				return err
			}

			var message string
			for i := int64(1); i <= parcount; i++ {
				key := strconv.FormatInt(i, 10)
				message += parts[key].Message
			}
			p.Message = message
			p.MessageID = parts["1"].ID

			return nil
		}
		return ErrInsufficientParts
	}
	return nil
}

func (s *SMSService) ProcessMO(p types.ProcessMOParams, r *types.NoReply) error {
	ctx := context.Background()

	// Check for multipart
	err := s.checkMultiPart(&p)
	if err != nil {
		if err == ErrInsufficientParts {
			return nil
		}
		return err
	}

	// check if we have a reply to a send sms
	replyNotFound := false
	sms, err := s.db.FindSMSRelatedToMO(ctx, p.To, p.From)
	if err != nil {
		if !errors.As(err, &errorlib.NotFoundErr{}) {
			log.Printf("[Processing MO] Error searching for related sms (To[%s] From[%s]): %s", p.To, p.From, err)
			return err
		}

		replyNotFound = true
	}

	var accountID string
	if replyNotFound {
		replySenders, err := s.senderRPC.FindSendersByAddress(ctx, &senderpb.FindSendersByAddressParams{
			Address: p.To,
		})
		if err != nil {
			log.Printf("[Processing MO] Could not find account for Sender: %s, error: %s", p.To, err)
			return err
		}

		nbSenders := len(replySenders.GetSenders())

		// check if we have 0 or multiple accounts related to this sender
		if nbSenders != 1 {
			msg := fmt.Sprintf("found %d account(s) for Sender: %s", nbSenders, p.To)
			log.Printf("[Processing MO] Error %s", msg)
			return errors.New(msg)
		}

		accountID = replySenders.GetSenders()[0].GetAccountId()
	} else {
		accountID = sms.AccountID
	}

	log.Printf("[Processing MO] Found sms (To[%s] From[%s]), accountID: %s", p.To, p.From, accountID)

	// Let the optout service deal with it if its an optout
	if err := s.optOutRPC.OptOutViaMsg(optOutRPC.OptOutViaMsgParams{
		AccountID:   accountID,
		Message:     p.Message,
		MessageType: `sms`,
		MessageID:   sms.ID,
	}); err != nil {
		log.Printf("[Processing MO] Error checking for OptOut: %s", err.Error())
		return err
	}

	var lastMessage *webhookpb.Message
	if !replyNotFound {
		lastMessage = &webhookpb.Message{
			Type:       "sms",
			Id:         sms.ID,
			Recipient:  sms.Recipient,
			Sender:     sms.Sender,
			Message:    sms.Message,
			MessageRef: sms.MessageRef,
		}
	} else {
		lastMessage = nil
	}

	_, err = s.webhookRPC.PublishMO(ctx, &webhookpb.PublishMOParams{
		AccountId:   accountID,
		SMSId:       p.MessageID,
		Recipient:   p.To,
		Sender:      p.From,
		Message:     p.Message,
		ReceivedAt:  timestamppb.Now(),
		LastMessage: lastMessage,
	})

	return err
}
