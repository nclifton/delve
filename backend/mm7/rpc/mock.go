package rpc

type mockS3 struct {
	err error
}

func (m mockS3) PutS3Content(content []byte, bucket, key string) error {
	return m.err
}
