### POST /mms/message

Request
```
type MMSPOSTRequest struct {
	// -- rename -- Body         string   `json:"body" valid:"required"`
	// -- rename -- ResourceURLs []string `json:"resource_urls" valid:"required"`
	ContactRef   string   `json:"contact_ref"`
	Recipient    string   `json:"recipient"`
	Sender       string   `json:"sender" valid:"required"`
	Subject      string   `json:"subject"`
	
	// new additions
	Message 			string 	`json:"message"`
	Country 			string 	`json:"country"`
	ShortenURLs 	bool 		`json:"shorten_urls"`
	ContentURLs []string `json:"content_urls" valid:"required"`
	
	// facilitates matching messages on teh global webhooks
	// a new field to add on the MMS doc
	MessageRef	string	`json:"message_ref"`
}

```

Response
```
type MMSPOSTResponse struct {
	MMS MMS `json:"mms"`
}

type ValidationError struct {}
```

### POST /mms/bulk

Request
```
type MMSBulkPOSTRequest struct {
	Messages []MMSPOSTRequest
}
```

Response
```
type MMSBulkPOSTResponse struct {
	SentCount 	int 	`json:"sent_count"`
	Sent 				[]MMS `json:"mms"`
	ErrorCount 	int 	`json:"error_count"`
	Errors 			[]ValidationError `json:"errors"`
}
```


