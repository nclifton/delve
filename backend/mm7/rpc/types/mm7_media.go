package types

type MediaStoreParams struct {
	FileName    string
	ProviderKey string
	Extension   string
	Data        []byte
}
type MediaStoreReply struct {
	URL string
}
