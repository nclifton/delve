package rpc

import mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"

type mockS3 struct {
	err error
}

func (m mockS3) PutS3Content(content []byte, bucket, key string) error {
	return m.err
}

type mockMMS struct {
	err error
}

func (m mockMMS) UpdateStatus(p mms.UpdateStatusParams) (err error) {
	return m.err
}
