package rpc

import (
	"fmt"
	"testing"
)

func TestStore(t *testing.T) {

	testErr := fmt.Errorf("testerror")

	tests := []struct {
		name        string
		params      MM7MediaStoreParams
		s3          mockS3
		expectedURL string
		expectedErr error
	}{
		{
			name: "test happy path",
			params: MM7MediaStoreParams{
				FileName:    "123",
				ProviderKey: "fake",
				Extension:   ".png",
			},
			s3:          mockS3{},
			expectedURL: "https://s3.ap-southeast-2.amazonaws.com/mms.media/fake_123.png",
			expectedErr: nil,
		},
		{
			name: "test with s3 error",
			params: MM7MediaStoreParams{
				FileName:    "123",
				ProviderKey: "fake",
				Extension:   ".png",
			},
			s3: mockS3{
				err: testErr,
			},
			expectedURL: "",
			expectedErr: testErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm7 := MM7{
				db: &db{
					s3: test.s3,
				},
				configVar: ConfigVar{
					AWSRegion:      "ap-southeast-2",
					MMSMediaBucket: "mms.media",
				},
			}

			r := &MM7MediaStoreReply{}
			err := mm7.Store(test.params, r)
			if err != test.expectedErr {
				t.Errorf("unexpected error %+v", err)
			}

			if r.URL != test.expectedURL {
				t.Errorf("expected URL %s, \nbut got %s", test.expectedURL, r.URL)
			}
		})

	}
}
