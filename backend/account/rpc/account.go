package rpc

import (
	"time"

	"github.com/burstsms/mtmo-tp/backend/account/rpc/types"
)

func (s *AccountService) FindByAPIKey(p types.FindByAPIKeyParams, r *types.FindByAPIKeyReply) error {

	var account *types.Account

	err := s.db.redis.Cached(
		"Account.FindByAPIKey:"+p.Key,
		&account,
		time.Minute*5,
		func() (interface{}, error) {
			return s.db.FindByAPIKey(p.Key)
		},
	)
	if err != nil {
		return err
	}
	r.Account = account
	return nil
}

func (s *AccountService) FindBySender(p types.FindBySenderParams, r *types.FindBySenderReply) error {
	var account *types.Account

	err := s.db.redis.Cached(
		"Account.FindBySender:"+p.Sender,
		&account,
		time.Minute*5,
		func() (interface{}, error) {
			return s.db.FindBySenderSMS(p.Sender)
		},
	)
	if err != nil {
		return err
	}

	r.Account = account
	return nil
}
