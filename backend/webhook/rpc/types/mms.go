package types

import (
	"time"
)

type PublishMMSStatusUpdateParams struct {
	AccountID         string    `json:"account_id"`
	MMSID             string    `json:"mms_id"`
	MessageRef        string    `json:"message_ref"`
	Recipient         string    `json:"recipient"`
	Sender            string    `json:"sender"`
	Status            string    `json:"status"`
	StatusDescription string    `json:"status_description"`
	StatusUpdatedAt   time.Time `json:"status_updated_at"`
}
