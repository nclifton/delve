package rpc

import (
	"fmt"

	"github.com/burstsms/mtmo-tp/backend/mm7/rpc/types"
)

func (s *MM7) Store(p types.MM7MediaStoreParams, r *types.MM7MediaStoreReply) error {
	key := fmt.Sprintf("%s_%s%s", p.ProviderKey, p.FileName, p.Extension)

	if err := s.svc.S3.PutS3Content(p.Data, s.configVar.MMSMediaBucket, key); err != nil {
		return err
	}

	// TODO: return key instead of the url?
	r.URL = fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", s.configVar.AWSRegion, s.configVar.MMSMediaBucket, key)

	return nil
}
