package rpc

import (
	"time"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc"
)

type OptOut struct {
	ID          string
	AccountID   string
	AccountName string
	Sender      string
	LinkID      string
	SMSID       string
	MMSID       string
}

type FindByLinkIDParams struct {
	LinkID string
}

type FindByLinkIDReply struct {
	*OptOut
}

func (s *OptOutService) FindByLinkID(p FindByLinkIDParams, r *FindByLinkIDReply) error {
	r.OptOut = &OptOut{
		ID:          `FakeOptoutID`,
		AccountName: `FakeAccountName`,
		AccountID:   `9b08870c-c6e1-461d-a06f-3f0078fde4fe`,
		Sender:      `61455678909`,
		LinkID:      p.LinkID,
	}
	return nil
}

type OptOutViaLinkParams struct {
	LinkID string
}

type OptOutViaLinkReply struct {
	*OptOut
}

func (s *OptOutService) OptOutViaLink(p OptOutViaLinkParams, r *OptOutViaLinkReply) error {

	var err error
	var optout FindByLinkIDReply

	err = s.FindByLinkID(FindByLinkIDParams{LinkID: p.LinkID}, &optout)
	if err != nil {
		return err
	}

	/*	if optout.SMSID != "" {
		} else if optout.OptOut.MMSID != "" {
		}*/

	err = s.webhookRPC.PublishOptOut(rpc.PublishOptOutParams{
		Source:    "link_hit",
		Timestamp: time.Now().UTC(),
		AccountID: optout.AccountID,
	})
	if err != nil {
		return err
	}

	r.OptOut = optout.OptOut

	return nil
}
