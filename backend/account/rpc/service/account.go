package service

import (
	"context"
	"time"

	"github.com/burstsms/mtmo-tp/backend/account/rpc/accountpb"
	"github.com/burstsms/mtmo-tp/backend/account/rpc/db"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *accountImpl) FindAccountByAPIKey(ctx context.Context, p *accountpb.FindAccountByAPIKeyParams) (*accountpb.FindAccountByAPIKeyReply, error) {
	accountDB := db.Account{}

	if err := a.redis.Cached(
		"Account.FindByAPIKey:"+p.GetKey(),
		&accountDB,
		time.Minute*5,
		func() (interface{}, error) {
			return a.db.FindAccountByAPIKey(ctx, p.GetKey())
		},
	); err != nil {
		return nil, err
	}

	return &accountpb.FindAccountByAPIKeyReply{
		Account: dbAccountToAccountPB(accountDB),
	}, nil
}

func (a *accountImpl) FindAccountByID(ctx context.Context, p *accountpb.FindAccountByIDParams) (*accountpb.FindAccountByIDReply, error) {
	accountDB := db.Account{}

	if err := a.redis.Cached(
		"Account.FindByID:"+p.GetId(),
		&accountDB,
		time.Minute*5,
		func() (interface{}, error) {
			return a.db.FindAccountByID(ctx, p.GetId())
		},
	); err != nil {
		return nil, err
	}

	return &accountpb.FindAccountByIDReply{
		Account: dbAccountToAccountPB(accountDB),
	}, nil
}

func dbAccountToAccountPB(a db.Account) *accountpb.Account {
	return &accountpb.Account{
		Id:             a.ID,
		Name:           a.Name,
		AlarisUsername: a.AlarisUsername,
		AlarisPassword: a.AlarisPassword,
		AlarisUrl:      a.AlarisURL,
		CreatedAt:      timestamppb.New(a.CreatedAt),
		UpdatedAt:      timestamppb.New(a.UpdatedAt),
	}
}
