package handler

import (
	"context"
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
)

type handler struct {
	log *logger.StandardLogger
}

func New() *handler {
	return &handler{
		log: logger.NewLogger(),
	}
}

func (h *handler) OnFinalFailure(ctx context.Context, body []byte) error {
	return nil
}

func (h *handler) Handle(ctx context.Context, body []byte, headers map[string]interface{}) error {
	log.Println("Test MgageSubmitHandler...")
	return nil
}
