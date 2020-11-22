package biz_test

import (
	"testing"

	"github.com/burstsms/mtmo-tp/backend/sms/biz"
)

const (
	validDLRText     = "id:5e4cc2e2d2a1a64e039a1bcb sub:001 delivrd:1 submit date:202008190505 done date:202008190505 stat:DELIVRD err:0 text:Bob creates a world which is good and fit for mankind"
	invalidDLRText   = "I am not a valid dlr"
	unknownDLRStatus = "id:5e4cc2e2d2a1a64e039a1bcb sub:001 delivrd:1 submit date:202008190505 done date:202008190505 stat:IAMNOTAVALIDDLRSTATUS err:0 text:Bob creates a world which is good and fit for mankind"
)

func TestParseDLRStatus(t *testing.T) {
	t.Run("valid dlr returns expected status", func(t *testing.T) {
		result, err := biz.ParseDLRStatus(validDLRText)
		if err != nil {
			t.Fatal("unexpected failure getting status from valid dlr:", err)
		}
		if result != "" {
			t.Fatalf("unexpected status returned from valid dlr, %s", result)
		}
	})

	t.Run("unknown dlr status is still returned", func(t *testing.T) {
		result, err := biz.ParseDLRStatus(unknownDLRStatus)
		if err != nil {
			t.Fatal("unexpected failure getting unknown status from valid dlr:", err)
		}
		if result != "IAMNOTAVALIDDLRSTATUS" {
			t.Fatal("unexpected status returned from valid dlr with unknown status")
		}
	})

	t.Run("return err if dlr is not valid", func(t *testing.T) {
		_, err := biz.ParseDLRStatus(invalidDLRText)
		if err == nil {
			t.Fatal("expected to get an error for invalid dlr")
		}
	})
}
