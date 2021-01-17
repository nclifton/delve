package errorlib

import "errors"

// TODO: refactor errors into generic subtypes

// sms, mms
var (
	ErrInvalidMobileNumber   = errors.New("Invalid mobile number")
	ErrInvalidPhoneNumber    = errors.New("Invalid mobile number")
	ErrInvalidSender         = errors.New("Invalid sender")
	ErrInvalidSenderNotFound = errors.New("Invalid sender: not found")
	ErrInvalidSenderCountry  = errors.New("Invalid sender: incorrect country")
	ErrInvalidSenderChannel  = errors.New("Invalid sender: incorrect channel")
	ErrInvalidSenderAddress  = errors.New("Invalid sender: incorrect address")

	ErrInvalidSMSTooManyParts = errors.New("SMS is too long")
	ErrInsufficientBalance    = errors.New("Insufficient balance")
	ErrNoKannelKey            = errors.New("No route key found for route")
	ErrInvalidRoute           = errors.New("No valid route for recipient")
	ErrSendSMS                = errors.New("Failed to send SMS")
	ErrProcessWebhooks        = errors.New("Failed to process webhooks")
)
