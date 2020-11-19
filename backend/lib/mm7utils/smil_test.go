package mm7utils

import (
	"fmt"
	"strings"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
)

func testContains(haystack string, needle string) cmp.Comparison {
	return func() cmp.Result {
		if strings.Contains(haystack, needle) {
			return cmp.ResultSuccess
		}
		return cmp.ResultFailure(fmt.Sprintf(`%s did not contain %s`, haystack, needle))
	}
}

func TestRenderSMIL(t *testing.T) {

	tests := []struct {
		attachments []SMILMedia
	}{{
		attachments: []SMILMedia{{ContentID: "<msg-txt>", MediaType: "txt"}, {ContentID: "<image1>", MediaType: "img"}},
	}}
	for _, testdata := range tests {
		smildata, err := renderSMIL(testdata.attachments)
		if err != nil {
			t.Error(err)
		}

		for idx, attachment := range testdata.attachments {
			assert.Assert(t, testContains(string(smildata), fmt.Sprintf(`<%s src="cid:%s" region="%s-%d"/>`, attachment.MediaType, attachment.ContentID, attachment.MediaType, idx)))
			assert.Assert(t, testContains(string(smildata), fmt.Sprintf(`<region id="%s-%d" top="50%%" left="0" height="50%%" width="100%%" fit="hidden"/>`, attachment.MediaType, idx)))
		}
	}

}
