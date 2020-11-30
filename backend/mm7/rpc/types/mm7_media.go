package types

type MM7MediaStoreParams struct {
	FileName    string
	ProviderKey string
	Extension   string
	Data        []byte
}
type MM7MediaStoreReply struct {
	URL string
}
