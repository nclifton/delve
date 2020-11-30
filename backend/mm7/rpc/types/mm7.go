package types

type PingResponse struct {
	Res string
}

type MM7SendParams struct {
	ID          string
	Subject     string
	Message     string
	Sender      string
	Recipient   string
	ContentURLs []string
	ProviderKey string
}

type MM7ProviderSpecParams struct {
	ProviderKey string
}

type MM7ProviderSpecReply struct {
	ProviderKey    string
	ImageSizeMaxKB int
}

type MM7UpdateStatusParams struct {
	ID          string
	MessageID   string
	Status      string
	Description string
}

type MM7DLRParams struct {
	ID          string
	Status      string
	Description string
}

type MM7DeliverParams struct {
	Subject     string
	Message     string
	Sender      string
	Recipient   string
	ContentURLs []string
	ProviderKey string
}

type MM7GetCachedContentParams struct {
	ContentURL string
}

type MM7GetCachedContentReply struct {
	Content []byte
}

type MM7CheckRateLimitParams struct {
	ProviderKey string
}

type MM7CheckRateLimitReply struct {
	Allow bool
}
