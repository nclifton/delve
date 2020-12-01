package types

import "time"

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

type FindTrackLinkByTrackLinkIDParams struct {
	TrackLinkID string
}

type FindTrackLinkByTrackLinkIDReply struct {
	TrackLink *TrackLink
}

type LinkHitParams struct {
	TrackLinkID string
}

type LinkHitReply struct {
	TrackLink *TrackLink
}
