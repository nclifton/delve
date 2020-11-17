package auth

import (
	"log"
	"net/http"

	account "github.com/burstsms/mtmo-tp/backend/account/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/api/middleware/context"
)

// values we set in the "authed_by" context key
const (
	AuthedByAPIKey = "api_key"
)

type auth struct {
	handler http.Handler
	opts    *Options
}

// Options lets us configure the behaviour of the middleware
type Options struct {
	UseAPIKey           bool
	FindAccountByAPIKey func(key string) (*account.Account, error)
}

func (a *auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// these are assumed to have values in a few places
	// panics happen if they aren't typed so don't remove ;p
	// TODO completely remove using context to store values in a global map
	context.Set(r, "authed_by", "")

	// authed means good to continue on the chain
	authed := false

	if !authed && a.opts.UseAPIKey {
		// check that api key has been provided
		key := r.Header.Get("x-api-key")

		if key != "" {
			// verify the token by fetching account via api key
			account, err := a.opts.FindAccountByAPIKey(key)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			context.Set(r, "authed_by", AuthedByAPIKey)
			context.Set(r, "account", account)
			authed = true
		}
	}

	if !authed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	a.handler.ServeHTTP(w, r)
}

// New auth middleware
func New(opts *Options) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &auth{
			handler: h,
			opts:    opts,
		}
	}
}
