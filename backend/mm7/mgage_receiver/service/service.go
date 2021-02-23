package service

import (
	"log"
	"net/http"

	"github.com/burstsms/mtmo-tp/backend/lib/rest"
	"github.com/julienschmidt/httprouter"
)

type Service struct{}

func Routes(router *httprouter.Router, baseRoute, authRoute rest.Handle) {
	s := Service{}

	router.NotFound = rest.HTTPHandlerWrapper(baseRoute(NotFoundRoute))
	router.PanicHandler = PanicRoute

	router.Handle("OPTIONS", "/v1", baseRoute(EmptyRoute))
	router.Handle("OPTIONS", "/v1/*path", baseRoute(EmptyRoute))

	router.POST("/v1/mms/dlr", authRoute(s.DLRPOST))
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
