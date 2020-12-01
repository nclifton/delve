package types

type SendParams struct {
	ID          string
	Subject     string
	Message     string
	Sender      string
	Recipient   string
	ContentURLs []string
	ProviderKey string
}

type ProviderSpecParams struct {
	ProviderKey string
}

type ProviderSpecReply struct {
	ProviderKey    string
	ImageSizeMaxKB int
}

type UpdateStatusParams struct {
	ID          string
	MessageID   string
	Status      string
	Description string
}

type DLRParams struct {
	ID          string
	Status      string
	Description string
}

type DeliverParams struct {
	Subject     string
	Message     string
	Sender      string
	Recipient   string
	ContentURLs []string
	ProviderKey string
}

type GetCachedContentParams struct {
	ContentURL string
}

type GetCachedContentReply struct {
	Content []byte
}

type CheckRateLimitParams struct {
	ProviderKey string
}

type CheckRateLimitReply struct {
	Allow bool
}
