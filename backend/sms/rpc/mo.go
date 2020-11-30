package rpc

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	optOutRPC "github.com/burstsms/mtmo-tp/backend/optout/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/sms/worker/msg"
	webhookRPC "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
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

var ErrInsufficientParts = errors.New("Insuffcient parts to combine MO")

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

		// do we habe all the parts
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

	// Check for multipart
	err := s.checkMultiPart(&p)
	if err != nil {
		if err == ErrInsufficientParts {
			return nil
		}
		return err
	}

	// Find the account from the sender
	account, err := s.accountRPC.FindBySender(p.To)
	if err != nil {
		log.Printf("[Processing MO] Could not find account for Sender: %s %s", p.To, err)
		return err
	}

	// check if we are a reply to a send sms
	sms, err := s.db.FindSMSRelatedToMO(account.Account.ID, p.From, p.To)
	if err != nil {
		log.Printf("[Processing MO] Error searching for related sms: %s %T", err.Error(), err)
		return err
	}

	log.Printf("[Processing MO] Found Account: %+v Related to: %+v", account, sms)

	// Let the optout service deal with it if its an optout
	err = s.optOutRPC.OptOutViaMsg(optOutRPC.OptOutViaMsgParams{
		AccountID:   account.Account.ID,
		Message:     p.Message,
		MessageType: `sms`,
		MessageID:   sms.ID,
	})
	if err != nil {
		log.Printf("[Processing MO] Error checking for OptOut: %s", err.Error())
		return err
	}

	var lastMessage *webhookRPC.LastMessage
	if *sms != (types.SMS{}) {
		lastMessage = &webhookRPC.LastMessage{
			Type:       "sms",
			ID:         sms.ID,
			Recipient:  sms.Recipient,
			Sender:     sms.Sender,
			Message:    sms.Message,
			MessageRef: sms.MessageRef,
		}
	} else {
		lastMessage = nil
	}

	return s.webhookRPC.PublishMO(webhookRPC.PublishMOParams{
		AccountID:   account.Account.ID,
		SMSID:       p.MessageID,
		Recipient:   p.To,
		Sender:      p.From,
		Message:     p.Message,
		ReceivedAt:  time.Now(),
		LastMessage: lastMessage,
	})
}
