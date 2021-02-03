package service

import (
	"errors"
	"log"
	"net/http"

	account "github.com/burstsms/mtmo-tp/backend/account/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/lib/rest"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"github.com/julienschmidt/httprouter"
)

const ClientsKey = "clients"

type Clients struct {
	SenderClient  senderpb.ServiceClient
	WebhookClient webhookpb.ServiceClient
	AccountClient *account.Client
	SMSClient     *sms.Client
	MMSClient     *mms.Client
}

type Service struct {
	*Clients
}

func accountFromCtx(hc *rest.HandlerContext) *account.Account {
	auth := hc.FromContext(rest.AuthKey)
	if auth == nil {
		hc.LogFatal(errors.New("accountFromCtx: Could not retrieve 'auth' from ctx"))
	}

	account, ok := auth.(*account.Account)
	if !ok {
		hc.LogFatal(errors.New("accountFromCtx: Could not coerce 'auth' to type *account.Account"))
	}

	return account
}

func Routes(router *httprouter.Router, baseRoute, authRoute rest.Handle, clients *Clients) {
	s := Service{clients}

	router.NotFound = rest.HTTPHandlerWrapper(baseRoute(NotFoundRoute))
	router.PanicHandler = PanicRoute

	router.Handle("OPTIONS", "/v1", baseRoute(EmptyRoute))
	router.Handle("OPTIONS", "/v1/*path", baseRoute(EmptyRoute))

	router.GET("/v1/sender", authRoute(s.SenderListGET))

	router.POST("/v1/sms", authRoute(s.SMSPOST))

	router.POST("/v1/mms", authRoute(s.MMSPOST))

	router.POST("/v1/webhook", authRoute(s.WebhookCreatePOST))
	router.PUT("/v1/webhook/:id", authRoute(s.WebhookUpdatePUT))
	router.GET("/v1/webhook/:id", authRoute(s.WebhookGET))
	router.DELETE("/v1/webhook/:id", authRoute(s.WebhookDELETE))
	router.GET("/v1/webhook", authRoute(s.WebhookListGET))
}

func EmptyRoute(_ *rest.HandlerContext) {}

func NotFoundRoute(hc *rest.HandlerContext) {
	hc.WriteJSONError("Not found", http.StatusNotFound, nil)
}

// Fallback required for panics and use of hc.LogFatal()
func PanicRoute(w http.ResponseWriter, _ *http.Request, err interface{}) {
	log.Printf("panic: %s", err)
	w.WriteHeader(500)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(`{ "error": "Internal server error" }`))
	if err != nil {
		log.Print(err)
	}
}
