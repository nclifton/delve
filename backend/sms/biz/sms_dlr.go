package biz

import (
	"errors"
	"regexp"
)

var (
	intermediateDLR = regexp.MustCompile("^(ACK|NACK)/(?:([0-9bx]+))?/?(?:(.*))?$")
	handsetDLR      = regexp.MustCompile(`^id:(.*) sub:(.*) de?li?vrd:(.*) submit date:(.*) done date:(.*) stat:(.*) err:(.*) (?i)(t)ext:(.*)`)
)

const (
	dlrCodeDelivered     = "DELIVRD"
	dlrCodeAccepted      = "ACCEPTD"
	dlrCodeExpired       = "EXPIRED"
	dlrCodeDeleted       = "DELETED"
	dlrCodeUndelivered   = "UNDELIV"
	dlrCodeRejected      = "REJECTD"
	dlrCodeSystemExpired = "SYSTEM_EXPIRED"
	dlrCodeAck           = "ACK"
	dlrCodeNack          = "NACK"
)

// DLR validation errors
var (
	ErrUnknownDLRFormat = errors.New("passed DLR does not match a known format")
)

/*
	refer to: https://github.com/burstsms/burstsms-system/src/681202bd52d55d52f4bab263c15aeca1ed1f9e72/legacy/src/Legacy/smppHandler.php#lines-919
	for legacy implementation
*/

// ParseDLR returns the status for the given dlr string
func ParseDLRStatus(dlrText string) (string, error) {
	var status string
	var statusTxt string

	if matches := intermediateDLR.FindStringSubmatch(dlrText); len(matches) >= 2 {
		statusTxt = matches[1]
	} else if matches := handsetDLR.FindStringSubmatch(dlrText); len(matches) >= 7 {
		statusTxt = matches[6]
	} else {
		return "", ErrUnknownDLRFormat
	}

	switch statusTxt {
	case dlrCodeDelivered:
		status = "" //db.DLRStatusDelivered
	case dlrCodeExpired:
		fallthrough
	case dlrCodeDeleted:
		fallthrough
	case dlrCodeRejected:
		fallthrough
	case dlrCodeSystemExpired:
		fallthrough
	case dlrCodeNack:
		status = "" //db.DLRStatusSoftBounce
	case dlrCodeUndelivered:
		status = "" //db.DLRStatusHardBounce
	case dlrCodeAccepted:
		fallthrough
	case dlrCodeAck:
		status = "" //db.DLRStatusAccepted
	default:
		status = statusTxt
	}

	return status, nil
}
