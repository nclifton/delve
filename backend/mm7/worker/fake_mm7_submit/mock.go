package fakemm7submitworker

import (
	"text/template"

	tcl "github.com/burstsms/mtmo-tp/backend/lib/tecloo/client"
	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc"
)

type mockRPCClient struct {
	rateLimitReply     mm7RPC.MM7CheckRateLimitReply
	cachedContentReply mm7RPC.MM7GetCachedContentReply
	err                error
}

func (m mockRPCClient) UpdateStatus(p mm7RPC.MM7UpdateStatusParams) (r *mm7RPC.NoReply, err error) {
	return nil, m.err
}

func (m mockRPCClient) CheckRateLimit(p mm7RPC.MM7CheckRateLimitParams) (r *mm7RPC.MM7CheckRateLimitReply, err error) {
	return &m.rateLimitReply, m.err
}

func (m mockRPCClient) GetCachedContent(p mm7RPC.MM7GetCachedContentParams) (r *mm7RPC.MM7GetCachedContentReply, err error) {
	return &m.cachedContentReply, m.err
}

type mockTecloo struct {
	postMM7Response tcl.PostMM7Response
	statusCode      int
	err             error
}

func (m mockTecloo) PostMM7(params tcl.PostMM7Params, soaptmpl *template.Template) (tcl.PostMM7Response, int, error) {
	return m.postMM7Response, m.statusCode, m.err
}
