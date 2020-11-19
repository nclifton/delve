package MgageSubmitworker

import (
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

func (h *MgageSubmitHandler) OnFinalFailure(body []byte) error {
	return nil
}

func (h *MgageSubmitHandler) Handle(body []byte, headers map[string]interface{}) error {
	log.Println("Test MgageSubmitHandler...")
	return nil
}
