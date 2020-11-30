package client

import (
	"github.com/burstsms/mtmo-tp/backend/track_link/rpc"
)

type GenerateTrackLinksParams = rpc.GenerateTrackLinksParams
type GenerateTrackLinksReply = rpc.GenerateTrackLinksReply

func (c *Client) GenerateTrackLinks(p GenerateTrackLinksParams) (r *GenerateTrackLinksReply, err error) {
	r = &GenerateTrackLinksReply{}
	err = c.Call("GenerateTrackLinks", p, r)
	return r, err
}

type FindTrackLinkByTrackLinkIDParams = rpc.FindTrackLinkByTrackLinkIDParams
type FindTrackLinkByTrackLinkIDReply = rpc.FindTrackLinkByTrackLinkIDReply

func (c *Client) FindTrackLinkByTrackLinkID(p FindTrackLinkByTrackLinkIDParams) (r *FindTrackLinkByTrackLinkIDReply, err error) {
	r = &FindTrackLinkByTrackLinkIDReply{}
	err = c.Call("FindTrackLinkByTrackLinkID", p, r)
	return r, err
}

type LinkHitParams = rpc.LinkHitParams
type LinkHitReply = rpc.LinkHitReply

func (c *Client) LinkHit(p LinkHitParams) (r *LinkHitReply, err error) {
	r = &LinkHitReply{}
	err = c.Call("IncrementTrackLinkHits", p, r)
	return r, err
}
