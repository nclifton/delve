### MMS Service -> MMS Send Queue

Queued Message
```
type Job struct {
	ID          objectid.ObjectID `json:"_id"`
	AccountID   objectid.ObjectID `json:"account_id"`
	CampaignID  objectid.ObjectID `json:"campaign_id"`
	ContactID   objectid.ObjectID `json:"contact_id"`
	Sender      string            `json:"sender"`
	Subject     string            `json:"subject"`
	// -- rename -- Body        string            `json:"body"`
	ContentURLs []string          `json:"content_urls"`
	Recipient   string            `json:"recipient"`
	ProviderKey string            `json:"provider_key"`
	
	// new
	Message				string 					`json:"message"`
}
```


