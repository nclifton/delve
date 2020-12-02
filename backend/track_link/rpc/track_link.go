package rpc

import (
	"context"
	"regexp"
	"strings"
	"time"

	mmstypes "github.com/burstsms/mtmo-tp/backend/mms/rpc/types"
	smstypes "github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/track_link/rpc/types"
	wtypes "github.com/burstsms/mtmo-tp/backend/webhook/rpc/types"
)

// URLRegex for detecting URLs in a string
var URLRegex = regexp.MustCompile(`(?:^|[\s]{1})(http(s?):\/\/)[^\s]+[^.;:(.\s)\s]`)

func (s *TrackLinkService) GenerateTrackLinks(p types.GenerateTrackLinksParams, r *types.GenerateTrackLinksReply) error {
	ctx := context.Background()

	msg := p.Message

	urls := URLRegex.FindAllString(p.Message, -1)
	for _, url := range urls {
		tracklink, err := s.db.InsertTrackLink(ctx, p.AccountID, p.MessageID, p.MessageType, url)
		if err != nil {
			return err
		}
		trackurl := "http://" + s.trackDomain + "/" + tracklink.TrackLinkID
		msg = strings.ReplaceAll(msg, url, trackurl)
	}

	r.Message = msg
	return nil
}

func (s *TrackLinkService) FindTrackLinkByTrackLinkID(p types.FindTrackLinkByTrackLinkIDParams, r *types.FindTrackLinkByTrackLinkIDReply) error {
	ctx := context.Background()

	tracklink, err := s.db.FindTrackLinkByTrackLinkID(ctx, p.TrackLinkID)
	if err != nil {
		return err
	}

	r.TrackLink = tracklink
	return nil
}

func (s *TrackLinkService) LinkHit(p types.LinkHitParams, r *types.LinkHitReply) error {
	ctx := context.Background()

	tracklink, err := s.db.IncrementTrackLinkHits(ctx, p.TrackLinkID)
	if err != nil {
		return err
	}

	r.TrackLink = tracklink

	var recipient, sender, message, messageref, subject string
	var contenturls []string
	// retrieve source msg
	switch tracklink.MessageType {
	case mmstypes.Name:
		msg, err := s.mmsRPC.FindByID(mmstypes.FindByIDParams{ID: tracklink.MessageID})
		if err != nil {
			return err
		}
		recipient = msg.MMS.Recipient
		sender = msg.MMS.Sender
		message = msg.MMS.Message
		messageref = msg.MMS.MessageRef
		subject = msg.MMS.Subject
		contenturls = msg.MMS.ContentURLs
	case smstypes.Name:
		msg, err := s.smsRPC.FindByID(smstypes.FindByIDParams{ID: tracklink.MessageID, AccountID: tracklink.AccountID})
		if err != nil {
			return err
		}
		recipient = msg.SMS.Recipient
		sender = msg.SMS.Sender
		message = msg.SMS.Message
		messageref = msg.SMS.MessageRef
	}

	sourcemsg := wtypes.SourceMessage{
		Type:        tracklink.MessageType,
		ID:          tracklink.MessageID,
		Recipient:   recipient,
		Sender:      sender,
		Message:     message,
		MessageRef:  messageref,
		Subject:     subject,
		ContentURLS: contenturls,
	}

	err = s.webhookRPC.PublishLinkHit(wtypes.PublishLinkHitParams{
		URL:           tracklink.URL,
		Hits:          tracklink.Hits,
		Timestamp:     time.Now().UTC(),
		SourceMessage: &sourcemsg,
		AccountID:     tracklink.AccountID,
	})
	return err
}
