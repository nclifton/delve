package rpc

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	mmsrpc "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
	smsrpc "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	wrpc "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
)

const optOutTemplate = "[opt-out-link]"

var optoutRegex = regexp.MustCompile(`\[opt-out-link\]`)

func (s *OptOutService) FindByLinkID(p types.FindByLinkIDParams, r *types.FindByLinkIDReply) error {

	ctx := context.Background()

	optOut, err := s.db.FindOptOutByLinkID(ctx, p.LinkID)
	if err != nil {
		return err
	}

	r.OptOut = optOut
	return nil
}

func (s *OptOutService) getOptOutOrigin(MessageType string, MessageID string, AccountID string) (*wrpc.SourceMessage, error) {

	var originMessage wrpc.SourceMessage
	// Get the linked message
	switch MessageType {
	case `sms`:
		sms, err := s.smsRPC.FindByID(smsrpc.FindByIDParams{ID: MessageID, AccountID: AccountID})
		if err != nil {
			return nil, nil
		}

		originMessage = wrpc.SourceMessage{
			Type:       `sms`,
			ID:         sms.ID,
			Recipient:  sms.Recipient,
			Sender:     sms.Sender,
			Message:    sms.Message,
			MessageRef: sms.MessageRef,
		}
	case `mms`:
		rply, err := s.mmsRPC.FindByID(mmsrpc.FindByIDParams{ID: MessageID})
		if err != nil {
			return nil, nil
		}

		originMessage = wrpc.SourceMessage{
			Type:        `mms`,
			ID:          rply.MMS.ID,
			Recipient:   rply.MMS.Recipient,
			Sender:      rply.MMS.Sender,
			Message:     rply.MMS.Message,
			MessageRef:  rply.MMS.MessageRef,
			Subject:     rply.MMS.Subject,
			ContentURLS: rply.MMS.ContentURLs,
		}

	default:
		return nil, fmt.Errorf("Invalid messageType (%s)", MessageType)
	}

	return &originMessage, nil
}

func (s *OptOutService) OptOutViaLink(p types.OptOutViaLinkParams, r *types.OptOutViaLinkReply) error {
	ctx := context.Background()

	optOut, err := s.db.FindOptOutByLinkID(ctx, p.LinkID)
	if err != nil {
		return err
	}

	originMessage, err := s.getOptOutOrigin(optOut.MessageType, optOut.MessageID, optOut.AccountID)
	if err != nil {
		return nil
	}

	if err := s.webhookRPC.PublishOptOut(wrpc.PublishOptOutParams{
		Source:        "link",
		Timestamp:     time.Now().UTC(),
		AccountID:     optOut.AccountID,
		OriginMessage: originMessage,
	}); err != nil {
		return err
	}

	r.OptOut = optOut
	return nil
}

func (s *OptOutService) GenerateOptOutLink(p types.GenerateOptOutLinkParams, r *types.GenerateOptOutLinkReply) error {
	ctx := context.Background()

	match := optoutRegex.FindAllString(p.Message, -1)
	if len(match) < 1 {
		r.Message = p.Message
		return nil
	}

	optOut, err := s.db.InsertOptOut(ctx, p.AccountID, p.MessageID, p.MessageType, p.Sender)
	if err != nil {
		return err
	}

	optOutURL := fmt.Sprintf("http://%s/%s", s.optOutDomain, optOut.LinkID)
	msg := strings.ReplaceAll(p.Message, optOutTemplate, optOutURL)

	r.Message = msg
	return nil
}
