package rpc

type SendParams struct {
	// -- rename -- Body         string   `json:"body" valid:"required"`
	// -- rename -- ResourceURLs []string `json:"resource_urls" valid:"required"`
	ContactRef string `json:"contact_ref"`
	Recipient  string `json:"recipient"`
	Sender     string `json:"sender" valid:"required"`
	Subject    string `json:"subject"`

	// new additions
	Message     string   `json:"message"`
	Country     string   `json:"country"`
	ShortenURLs bool     `json:"shorten_urls"`
	ContentURLs []string `json:"content_urls" valid:"required"`

	// facilitates matching messages on the global webhooks
	// a new field to add on the MMS doc
	MessageRef string `json:"message_ref"`
}

type SendReply struct {
}

func (s *MMSService) Send(p SendParams, r *SendReply) error {
	// TODO: to implement
	return nil
}
