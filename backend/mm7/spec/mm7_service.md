### RPC mm7.Send()

Request

```
type SendParams struct {
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
type ProviderSpecParams struct {
	ProviderKey	string
}
```

Reply

```
type ProviderSpecReply struct {
	ProviderKey			string
	ImageSizeMaxKB	int
}
```

### RPC mm7.UpdateStatus()

Request

```
type UpdateStatusParams struct {
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
type DLRParams struct {
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
type DeliverParams struct {
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
type GetCachedContentParams struct {
	ContentURL	string
}
```

Reply

```
type GetCachedContentReply struct {
	Content    []byte
}
```

### RPC mm7.CheckRateLimit()

Request

```
type CheckRateLimitParams struct {
	ProviderKey	string
}
```

Reply

```
type CheckRateLimitReply struct {
	Allow    bool
}
```
