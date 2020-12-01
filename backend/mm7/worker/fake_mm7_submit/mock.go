package fakemm7submitworker

import (
	"text/template"

	tcl "github.com/burstsms/mtmo-tp/backend/lib/tecloo/client"
	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"
)

type mockRPCClient struct {
	providerSpecReply  mm7RPC.ProviderSpecReply
	rateLimitReply     mm7RPC.CheckRateLimitReply
	cachedContentReply mm7RPC.GetCachedContentReply

	err error
}

func (m mockRPCClient) UpdateStatus(p mm7RPC.UpdateStatusParams) error {
	return m.err
}

func (m mockRPCClient) ProviderSpec(p mm7RPC.ProviderSpecParams) (r *mm7RPC.ProviderSpecReply, err error) {
	return &m.providerSpecReply, m.err
}

func (m mockRPCClient) CheckRateLimit(p mm7RPC.CheckRateLimitParams) (r *mm7RPC.CheckRateLimitReply, err error) {
	return &m.rateLimitReply, m.err
}

func (m mockRPCClient) GetCachedContent(p mm7RPC.GetCachedContentParams) (r *mm7RPC.GetCachedContentReply, err error) {
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
