package rpc

import (
	"strings"
	"time"

	"github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
	smsrpc "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	wrpc "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
)

var optOutWords = []string{"STOP", "END", "QUIT", "UN", "UNSUB", "UNSUBSCRIBE", "REMOVE", "RMV"}

func (s *OptOutService) OptOutViaMsg(p types.OptOutViaMsgParams, r *types.NoReply) error {

	var optOut bool

	// Does the message match the opt out keywords
	testMsg := strings.TrimSpace(p.Message)
	for _, kw := range optOutWords {
		if kw == strings.ToUpper(testMsg) {
			optOut = true
			break
		}
	}

	if optOut {
		var sourceMessage wrpc.SourceMessage

		// Get the linked message
		if p.MessageType == `sms` {
			sms, err := s.smsRPC.FindByID(smsrpc.FindByIDParams{ID: p.MessageID, AccountID: p.AccountID})
			if err != nil {
				return nil
			}

			sourceMessage = wrpc.SourceMessage{
				Type:       `sms`,
				ID:         sms.ID,
				Recipient:  sms.Recipient,
				Sender:     sms.Sender,
				Message:    sms.Message,
				MessageRef: sms.MessageRef,
			}

		}

		err := s.webhookRPC.PublishOptOut(wrpc.PublishOptOutParams{
			Source:        "sms_inbound",
			Timestamp:     time.Now().UTC(),
			AccountID:     p.AccountID,
			SourceMessage: &sourceMessage,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
