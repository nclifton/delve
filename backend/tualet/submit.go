package tualet

import (
	"log"
	"net/http"

	"github.com/burstsms/mtmo-tp/backend/logger"
)

func SubmitGET(r *Route) {

	params := r.r.URL.Query()
	username := params.Get("username")
	log.Printf("username: %s", username)
	if username == "" {
		r.WriteError("not authorized (check login and password)", http.StatusUnauthorized)
		return
	}
	password := params.Get("password")
	if password == "" {
		r.WriteError("not authorized (check login and password)", http.StatusUnauthorized)
		return
	}

	command := params.Get("command")
	if command != "submit" {
		r.WriteError("invalid command", http.StatusBadRequest)
		return
	}

	message := params.Get("message")
	if message == "" {
		r.WriteError("message missing", http.StatusBadRequest)
		return
	}

	dnis := params.Get("dnis")
	if message == "" {
		r.WriteError("destination missing", http.StatusBadRequest)
		return
	}

	ani := params.Get("ani")
	longMessageMode := params.Get("longMessageMode")

	r.api.log.Fields(logger.Fields{"msgid": "xxx", "dnis": dnis, "ani": ani, "message": message, "longMessageMode": longMessageMode, "status": 200})

	type payload struct {
		MessageID string `json:"message_id"`
	}
	data := payload{
		MessageID: "xxx",
	}

	r.Write(data, http.StatusOK)
}
