package service

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/db"
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

func (s *senderImpl) CreateSenders(ctx context.Context, r *senderpb.CreateSendersParams) (*senderpb.CreateSendersReply, error) {

	ss := []*senderpb.Sender{}

	if len(r.Senders) > 0 {

		newSenders := make([]db.Sender, 0, len(r.Senders))
		for _, sender := range r.Senders {
			newSenders = append(newSenders, db.Sender{
				AccountID:      sender.AccountId,
				Address:        sender.Address,
				MMSProviderKey: sender.MMSProviderKey,
				Channels:       sender.Channels,
				Country:        sender.Country,
				Comment:        sender.Comment,
			})
		}

		dbSenders, err := s.db.CreateSenders(ctx, newSenders)
		if err != nil {
			return nil, err
		}

		for _, s := range dbSenders {
			ss = append(ss, dbSenderToSender(s))
		}
	}

	return &senderpb.CreateSendersReply{
		Senders: ss,
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
