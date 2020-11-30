package client

import (
	wrpc "github.com/burstsms/mtmo-tp/backend/webhook/rpc/types"
)

type FindParams = wrpc.FindParams
type FindReply = wrpc.FindReply

type SourceMessage = wrpc.SourceMessage

func (c *Client) Find(accountID string) (r *FindReply, err error) {
	r = &FindReply{}
	err = c.Call("Find", FindParams{AccountID: accountID}, r)
	return
}

type InsertParams = wrpc.InsertParams
type InsertReply = wrpc.InsertReply

func (c *Client) Insert(accountID, event, name, url string, rateLimit int) (r *InsertReply, err error) {
	r = &InsertReply{}
	err = c.Call("Insert", InsertParams{
		AccountID: accountID,
		Event:     event,
		Name:      name,
		URL:       url,
		RateLimit: rateLimit,
	}, r)
	return
}

type DeleteParams = wrpc.DeleteParams
type DeleteReply = wrpc.NoReply

func (c *Client) Delete(accountID, id string) (r *DeleteReply, err error) {
	r = &DeleteReply{}
	err = c.Call("Delete", DeleteParams{
		AccountID: accountID,
		ID:        id,
	}, r)
	return
}

type UpdateParams = wrpc.UpdateParams
type UpdateReply = wrpc.UpdateReply

func (c *Client) Update(id, accountID, event, name, url string, rateLimit int) (r *UpdateReply, err error) {
	r = &UpdateReply{}
	err = c.Call("Update", UpdateParams{
		AccountID: accountID,
		ID:        id,
		Event:     event,
		Name:      name,
		URL:       url,
		RateLimit: rateLimit,
	}, r)
	return
}
