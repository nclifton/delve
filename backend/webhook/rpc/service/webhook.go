package service

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/db"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

func (s *webhookImpl) Insert(ctx context.Context, r *webhookpb.InsertParams) (*webhookpb.InsertReply, error) {
	dbWebhook, err := s.db.InsertWebhook(ctx, r.AccountId, r.Event, r.Name, r.URL, r.RateLimit)
	if err != nil {
		return nil, err
	}
	return &webhookpb.InsertReply{
		Webhook: dbWebhookToWebhook(dbWebhook),
	}, nil
}

// TODO requires integration test
func (s *webhookImpl) FindByID(ctx context.Context, r *webhookpb.FindByIDParams) (*webhookpb.FindByIDReply, error) {
	w, err := s.db.FindWebhookByID(ctx, r.AccountId, r.WebhookId)
	if err != nil {
		if errors.As(err, &errorlib.NotFoundErr{}) {
			err = status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return &webhookpb.FindByIDReply{
		Webhook: dbWebhookToWebhook(w),
	}, nil
}

func (s *webhookImpl) Find(ctx context.Context, r *webhookpb.FindParams) (*webhookpb.FindReply, error) {
	dbWebhooks, err := s.db.FindWebhook(ctx, r.AccountId)
	if err != nil {
		return nil, err
	}

	webhooks := []*webhookpb.Webhook{}
	for _, w := range dbWebhooks {
		webhooks = append(webhooks, dbWebhookToWebhook(w))
	}

	return &webhookpb.FindReply{
		Webhooks: webhooks,
	}, nil
}

func (s *webhookImpl) Delete(ctx context.Context, r *webhookpb.DeleteParams) (*webhookpb.NoReply, error) {
	err := s.db.DeleteWebhook(ctx, r.Id, r.AccountId)

	return &webhookpb.NoReply{}, err
}

func (s *webhookImpl) Update(ctx context.Context, r *webhookpb.UpdateParams) (*webhookpb.UpdateReply, error) {
	w, err := s.db.UpdateWebhook(ctx, r.Id, r.AccountId, r.Event, r.Name, r.URL, r.RateLimit)
	if err != nil {
		return nil, err
	}
	return &webhookpb.UpdateReply{
		Webhook: dbWebhookToWebhook(w),
	}, nil
}

func dbWebhookToWebhook(w db.Webhook) *webhookpb.Webhook {
	return &webhookpb.Webhook{
		Id:        w.ID,
		AccountId: w.AccountID,
		Event:     w.Event,
		Name:      w.Name,
		URL:       w.URL,
		RateLimit: w.RateLimit,
		CreatedAt: timestamppb.New(w.CreatedAt),
		UpdatedAt: timestamppb.New(w.UpdatedAt),
	}
}
