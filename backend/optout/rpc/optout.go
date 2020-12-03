package rpc

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
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

func (s *OptOutService) OptOutViaLink(p types.OptOutViaLinkParams, r *types.OptOutViaLinkReply) error {
	ctx := context.Background()

	optOut, err := s.db.FindOptOutByLinkID(ctx, p.LinkID)
	if err != nil {
		return err
	}

	if err := s.webhookRPC.PublishOptOut(wrpc.PublishOptOutParams{
		Source:    "link_hit",
		Timestamp: time.Now().UTC(),
		AccountID: optOut.AccountID,
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
