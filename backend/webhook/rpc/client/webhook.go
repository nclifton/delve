package client

import "github.com/burstsms/mtmo-tp/backend/webhook/rpc/types"

type FindParams = types.FindParams
type FindReply = types.FindReply

type SourceMessage = types.SourceMessage

func (c *Client) Find(accountID string) (r *FindReply, err error) {
	r = &FindReply{}
	err = c.Call("Find", FindParams{AccountID: accountID}, r)
	return
}

type InsertParams = types.InsertParams
type InsertReply = types.InsertReply

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

type DeleteParams = types.DeleteParams
type DeleteReply = types.NoReply

func (c *Client) Delete(accountID, id string) (r *DeleteReply, err error) {
	r = &DeleteReply{}
	err = c.Call("Delete", DeleteParams{
		AccountID: accountID,
		ID:        id,
	}, r)
	return
}

type UpdateParams = types.UpdateParams
type UpdateReply = types.UpdateReply

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
