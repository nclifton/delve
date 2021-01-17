package service

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/app/db"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

func (s *senderImpl) FindByAddress(ctx context.Context, r *senderpb.FindByAddressParams) (*senderpb.FindByAddressReply, error) {

	sender, err := s.db.SenderFindByAddress(ctx, r.AccountId, r.Address)
	if err != nil {
		if errors.As(err, &errorlib.NotFoundErr{}) {
			err = status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return &senderpb.FindByAddressReply{
		Sender: dbSenderToSender(sender),
	}, nil
}

func dbSenderToSender(sender db.Sender) *senderpb.Sender {
	return &senderpb.Sender{
		Id:             sender.ID,
		AccountId:      sender.AccountID,
		Address:        sender.Address,
		MMSProviderKey: sender.MMSProviderKey,
		Channels:       sender.Channels,
		Country:        sender.Country,
		Comment:        sender.Comment,
		CreatedAt:      timestamppb.New(sender.CreatedAt),
		UpdatedAt:      timestamppb.New(sender.UpdatedAt),
	}
}

func (s *senderImpl) FindByAccountId(ctx context.Context, r *senderpb.FindByAccountIdParams) (*senderpb.FindByAccountIdReply, error) {

	senders, err := s.db.SenderFindByAccountId(ctx, r.AccountId)
	if err != nil {
		return nil, err
	}
	ss := []*senderpb.Sender{}
	for _, s := range senders {
		ss = append(ss, dbSenderToSender(s))
	}

	return &senderpb.FindByAccountIdReply{
		Senders: ss,
	}, nil
}
