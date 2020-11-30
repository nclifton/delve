package rpc

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc"
)

const optOutTemplate = "[opt-out-link]"

var optoutRegex = regexp.MustCompile(`\[opt-out-link\]`)

type OptOut struct {
	ID          string
	AccountID   string
	MessageID   string
	MessageType string
	Sender      string
	LinkID      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FindByLinkIDParams struct {
	LinkID string
}

type FindByLinkIDReply struct {
	*OptOut
}

func (s *OptOutService) FindByLinkID(p FindByLinkIDParams, r *FindByLinkIDReply) error {
	ctx := context.Background()

	optOut, err := s.db.FindOptOutByLinkID(ctx, p.LinkID)
	if err != nil {
		return err
	}

	r.OptOut = optOut
	return nil
}

type OptOutViaLinkParams struct {
	LinkID string
}

type OptOutViaLinkReply struct {
	*OptOut
}

func (s *OptOutService) OptOutViaLink(p OptOutViaLinkParams, r *OptOutViaLinkReply) error {
	ctx := context.Background()

	optOut, err := s.db.FindOptOutByLinkID(ctx, p.LinkID)
	if err != nil {
		return err
	}

	if err := s.webhookRPC.PublishOptOut(rpc.PublishOptOutParams{
		Source:    "link_hit",
		Timestamp: time.Now().UTC(),
		AccountID: optOut.AccountID,
	}); err != nil {
		return err
	}

	r.OptOut = optOut
	return nil
}

type GenerateOptoutLinkParams struct {
	AccountID   string
	MessageID   string
	MessageType string
	Message     string
}

type GenerateOptoutLinkReply struct {
	Message string
}

func (s *OptOutService) GenerateOptoutLink(p GenerateOptoutLinkParams, r *GenerateOptoutLinkReply) error {
	ctx := context.Background()

	match := optoutRegex.FindAllString(p.Message, -1)
	if len(match) < 1 {
		r.Message = p.Message
		return nil
	}

	optOut, err := s.db.InsertOptOut(ctx, p.AccountID, p.MessageID, p.MessageType)
	if err != nil {
		return err
	}

	optOutURL := fmt.Sprintf("http://%s/%s", s.trackHost, optOut.LinkID)
	msg := strings.ReplaceAll(p.Message, optOutTemplate, optOutURL)

	r.Message = msg
	return nil
}
