package service

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

func (s *senderImpl) FindSenderByAddressAndAccountID(ctx context.Context, r *senderpb.FindSenderByAddressAndAccountIDParams) (*senderpb.FindSenderByAddressAndAccountIDReply, error) {
	sender, err := s.db.FindSenderByAddressAndAccountID(ctx, r.AccountId, r.Address)
	if err != nil {
		if errors.As(err, &errorlib.NotFoundErr{}) {
			err = status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return &senderpb.FindSenderByAddressAndAccountIDReply{
		Sender: dbSenderToSender(sender),
	}, nil
}

func (s *senderImpl) FindSendersByAccountId(ctx context.Context, r *senderpb.FindSendersByAccountIdParams) (*senderpb.FindSendersByAccountIdReply, error) {
	senders, err := s.db.FindSendersByAccountId(ctx, r.AccountId)
	if err != nil {
		return nil, err
	}

	ss := []*senderpb.Sender{}
	for _, s := range senders {
		ss = append(ss, dbSenderToSender(s))
	}

	return &senderpb.FindSendersByAccountIdReply{
		Senders: ss,
	}, nil
}

func (s *senderImpl) FindSendersByAddress(ctx context.Context, r *senderpb.FindSendersByAddressParams) (*senderpb.FindSendersByAddressReply, error) {
	senders, err := s.db.FindSendersByAddress(ctx, r.Address)
	if err != nil {
		return nil, err
	}

	ss := []*senderpb.Sender{}
	for _, s := range senders {
		ss = append(ss, dbSenderToSender(s))
	}

	return &senderpb.FindSendersByAddressReply{
		Senders: ss,
	}, nil
}
