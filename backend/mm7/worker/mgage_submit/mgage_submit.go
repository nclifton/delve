package MgageSubmitworker

import (
	"context"
	"log"
	"net/rpc"
)

type MgageSubmitHandler struct {
	mm7RPC *rpc.Client
}

func NewHandler(c *rpc.Client) *MgageSubmitHandler {
	return &MgageSubmitHandler{
		mm7RPC: c,
	}
}

func (h *MgageSubmitHandler) OnFinalFailure(ctx context.Context, body []byte) error {
	return nil
}

func (h *MgageSubmitHandler) Handle(ctx context.Context, body []byte, headers map[string]interface{}) error {
	log.Println("Test MgageSubmitHandler...")
	return nil
}
