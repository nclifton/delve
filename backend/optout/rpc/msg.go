package rpc

import (
	"strings"
	"time"

	"github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
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
		originMessage, err := s.getOptOutOrigin(p.MessageType, p.MessageID, p.AccountID)
		if err != nil {
			return nil
		}

		err = s.webhookRPC.PublishOptOut(wrpc.PublishOptOutParams{
			Source:        "sms",
			Timestamp:     time.Now().UTC(),
			AccountID:     p.AccountID,
			OriginMessage: originMessage,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
