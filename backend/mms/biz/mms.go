package biz

import (
	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	"github.com/burstsms/mtmo-tp/backend/lib/stringutil"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

func IsValidSender(sender *senderpb.Sender, address, country string) error {
	if sender == nil {
		return errorlib.ErrInvalidSenderNotFound
	}
	if sender.Address != address {
		return errorlib.ErrInvalidSenderAddress
	}
	if sender.Country != country {
		return errorlib.ErrInvalidSenderCountry
	}
	if !stringutil.Includes(sender.Channels, "mms") {
		return errorlib.ErrInvalidSenderChannel
	}
	return nil
}
