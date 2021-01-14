package service

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

func (s *senderImpl) FindByAddress(ctx context.Context, r *senderpb.FindByAddressParams) (*senderpb.FindByAddressReply, error) {

	data, err := s.db.SenderFindByAddress(ctx, r.AccountId, r.Address)
	if err != nil {
		if errors.As(err, &errorlib.NotFoundErr{}) {
			err = status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return &senderpb.FindByAddressReply{
		Sender: &senderpb.Sender{
			Id:             data.ID,
			AccountId:      data.AccountID,
			Address:        data.Address,
			MMSProviderKey: data.MMSProviderKey,
			Channels:       data.Channels,
			Country:        data.Country,
			Comment:        data.Comment,
			CreatedAt:      timestamppb.New(data.CreatedAt),
			UpdatedAt:      timestamppb.New(data.UpdatedAt),
		},
	}, nil
}
