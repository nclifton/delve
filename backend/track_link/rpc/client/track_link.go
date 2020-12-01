package client

import (
	"github.com/burstsms/mtmo-tp/backend/track_link/rpc/types"
)

type GenerateTrackLinksParams = types.GenerateTrackLinksParams
type GenerateTrackLinksReply = types.GenerateTrackLinksReply

func (c *Client) GenerateTrackLinks(p GenerateTrackLinksParams) (r *GenerateTrackLinksReply, err error) {
	r = &GenerateTrackLinksReply{}
	err = c.Call("GenerateTrackLinks", p, r)
	return r, err
}

type FindTrackLinkByTrackLinkIDParams = types.FindTrackLinkByTrackLinkIDParams
type FindTrackLinkByTrackLinkIDReply = types.FindTrackLinkByTrackLinkIDReply

func (c *Client) FindTrackLinkByTrackLinkID(p FindTrackLinkByTrackLinkIDParams) (r *FindTrackLinkByTrackLinkIDReply, err error) {
	r = &FindTrackLinkByTrackLinkIDReply{}
	err = c.Call("FindTrackLinkByTrackLinkID", p, r)
	return r, err
}

type LinkHitParams = types.LinkHitParams
type LinkHitReply = types.LinkHitReply

func (c *Client) LinkHit(p LinkHitParams) (r *LinkHitReply, err error) {
	r = &LinkHitReply{}
	err = c.Call("LinkHit", p, r)
	return r, err
}
