package rpc

import (
	"context"
	"regexp"
	"strings"
	"time"
)

// URLRegex for detecting URLs in a string
var URLRegex = regexp.MustCompile(`(?:^|[\s]{1})(http(s?):\/\/)[^\s]+[^.;:(.\s)\s]`)

type TrackLink struct {
	ID          string
	AccountID   string
	MessageID   string
	MessageType string
	TrackLinkID string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	URL         string
	Hits        int
}

type GenerateTrackLinksParams struct {
	AccountID   string
	MessageID   string
	MessageType string
	Message     string
}

type GenerateTrackLinksReply struct {
	Message string
}

func (s *TrackLinkService) GenerateTrackLinks(p GenerateTrackLinksParams, r *GenerateTrackLinksReply) error {
	ctx := context.Background()

	msg := p.Message

	urls := URLRegex.FindAllString(p.Message, -1)
	for _, url := range urls {
		tracklink, err := s.db.InsertTrackLink(ctx, p.AccountID, p.MessageID, p.MessageType, url)
		if err != nil {
			return err
		}
		trackurl := "http://" + s.trackHost + "/" + tracklink.TrackLinkID
		msg = strings.ReplaceAll(msg, url, trackurl)
	}

	r.Message = msg
	return nil
}

type FindTrackLinkByTrackLinkIDParams struct {
	AccountID   string
	TrackLinkID string
}

type FindTrackLinkByTrackLinkIDReply struct {
	TrackLink *TrackLink
}

func (s *TrackLinkService) FindTrackLinkByTrackLinkID(p FindTrackLinkByTrackLinkIDParams, r *FindTrackLinkByTrackLinkIDReply) error {
	ctx := context.Background()

	tracklink, err := s.db.FindTrackLinkByTrackLinkID(ctx, p.AccountID, p.TrackLinkID)
	if err != nil {
		return err
	}

	r.TrackLink = tracklink
	return nil
}

type LinkHitParams struct {
	AccountID   string
	TrackLinkID string
}

type LinkHitReply struct {
	TrackLink *TrackLink
}

func (s *TrackLinkService) LinkHit(p LinkHitParams, r *LinkHitReply) error {
	ctx := context.Background()

	tracklink, err := s.db.IncrementTrackLinkHits(ctx, p.AccountID, p.TrackLinkID)
	if err != nil {
		return err
	}

	r.TrackLink = tracklink
	return nil
}
