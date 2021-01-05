package rpc

import (
	"context"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
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

		// Get the linked message
		sourceMessage, err := s.getOptOutOrigin(p.MessageType, p.MessageID, p.AccountID)
		if err != nil {
			return err
		}
		_, err = s.webhookRPC.PublishOptOut(context.Background(), &webhookpb.PublishOptOutParams{
			Source:        "sms_inbound",
			Timestamp:     timestamppb.New(time.Now().UTC()),
			AccountId:     p.AccountID,
			SourceMessage: sourceMessage,
		})
		if err != nil {
			return err
		}

	}

	return nil
}
