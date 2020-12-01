package optussubmitworker

import (
	"text/template"

	optcl "github.com/burstsms/mtmo-tp/backend/lib/optus/client"
	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"
)

type mockRPCClient struct {
	rateLimitReply     mm7RPC.CheckRateLimitReply
	cachedContentReply mm7RPC.GetCachedContentReply
	err                error
}

func (m mockRPCClient) UpdateStatus(p mm7RPC.UpdateStatusParams) (err error) {
	return m.err
}

func (m mockRPCClient) CheckRateLimit(p mm7RPC.CheckRateLimitParams) (r *mm7RPC.CheckRateLimitReply, err error) {
	return &m.rateLimitReply, m.err
}

func (m mockRPCClient) GetCachedContent(p mm7RPC.GetCachedContentParams) (r *mm7RPC.GetCachedContentReply, err error) {
	return &m.cachedContentReply, m.err
}

type mockOptus struct {
	postMM7Response optcl.PostMM7Response
	statusCode      int
	err             error
}

func (m mockOptus) PostMM7(params optcl.PostMM7Params, soaptmpl *template.Template) (optcl.PostMM7Response, int, error) {
	return m.postMM7Response, m.statusCode, m.err
}
