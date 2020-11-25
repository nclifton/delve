### RPC mm7.Send()

Request

```
type MM7SendParams struct {
	ID					string
	Subject 		string
	Message 		string
	Sender			string
	Recipient		string
	ContentURLs	[]string
	ProviderKey	string
}
```

### RPC mm7.ProviderSpec()

Request

```
type MM7ProviderSpecParams struct {
	ProviderKey	string
}
```

Reply

```
type MM7ProviderSpecReply struct {
	ProviderKey			string
	ImageSizeMaxKB	int
}
```

### RPC mm7.UpdateStatus()

Request

```
type MM7UpdateStatusParams struct {
	MMSID       string
	MessageID   string
	Status      string
	Description string
}
```

Reply

```
No reply.
```

### RPC mm7.DLR()

Request

```
type MM7DLRParams struct {
	ID 					string
	Status			string
	Description	string
}
```

Reply

```
No reply.
```

### RPC mm7.Deliver()

Request

```
type MM7DeliverParams struct {
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

### RPC mm7.GetCachedContent()

Request

```
type MM7GetCachedContentParams struct {
	ContentURL	string
}
```

Reply

```
type MM7GetCachedContentReply struct {
	Content    []byte
}
```

### RPC mm7.CheckRateLimit()

Request

```
type MM7CheckRateLimitParams struct {
	ProviderKey	string
}
```

Reply

```
type MM7CheckRateLimitReply struct {
	Allow    bool
}
```
