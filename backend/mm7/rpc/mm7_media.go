package rpc

import (
	"fmt"
)

type MM7MediaStoreParams struct {
	FileName    string
	ProviderKey string
	Extension   string
	Data        []byte
}
type MM7MediaStoreReply struct {
	URL string
}

func (s *MM7) Store(p MM7MediaStoreParams, r *MM7MediaStoreReply) error {
	key := fmt.Sprintf("%s_%s%s", p.ProviderKey, p.FileName, p.Extension)

	if err := s.svc.S3.PutS3Content(p.Data, s.configVar.MMSMediaBucket, key); err != nil {
		return err
	}

	// TODO: return key instead of the url?
	r.URL = fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", s.configVar.AWSRegion, s.configVar.MMSMediaBucket, key)

	return nil
}
