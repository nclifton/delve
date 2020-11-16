### RPC mms.CreateAndQueue()
Request
```
type MMSCreateAndQueueParams struct {
	AccountID     db.OID
	Message       string
	Subject 			string
	// -- rename -- ResourceURLs 	[]string
	Recipient     string
	Sender        string
	Source        string   // TODO: figure out if this is required (events maybe?)
	Country       string
	// --delete-- EmbeddedContact  	db.EmbeddedContact
	
	// New
	ShortenURLs		bool
	ContentURLs 	[]string
	MessageRef		string
	ContactID			string
}
```

Reply
```
type MMSCreateAndQueueReply struct {
	MMS		*db.MMS
}
```

### RPC mms.UpdateStatus()

* Used for DLR too
* Put `status_updated_at` field on MMS doc

Request
```
type MMSUpdateStatusParams struct {
	ID          db.OID 	`valid:"objectid"`
	Status      string 	`valid:"required"`
	Description string	
}
```

Reply
```
No reply.
```

### RPC mms.Inbound()

* mm7 Deliver

Request
```
type MMSInboundParams struct {
	Subject 		string
	Message 		string
	Sender			string
	Recipient		string
	ContentURLs	[]string
	ProviderKey	string
}
```

Reply
```
No reply.
```


