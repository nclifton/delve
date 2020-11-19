package OptusSubmitworker

import (
	"log"
	"net/rpc"
)

type OptusSubmitHandler struct {
	mm7RPC *rpc.Client
}

func NewHandler(c *rpc.Client) *OptusSubmitHandler {
	return &OptusSubmitHandler{
		mm7RPC: c,
	}
}

func (h *OptusSubmitHandler) OnFinalFailure(body []byte) error {
	return nil
}

func (h *OptusSubmitHandler) Handle(body []byte, headers map[string]interface{}) error {
	log.Println("Test OptusSubmitHandler...")
	return nil
}
