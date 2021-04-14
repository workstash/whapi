package routes

import (
	"encoding/json"
	"net/http"

	"github.com/workstash/whapi/infrastructure/whats"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Structure to return JSON with Contacts Info

type ContactInfo struct {
	Exists string `json:"exists"`
	Status string `json:"status"`
	Online string `json:"online"`
	Thumb  string `json:"thumb"`
}

func checkContact() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validate data

		device, ok := r.URL.Query()["device"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequest)
			return
		}

		num, ok := r.URL.Query()["num"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequest)
			return
		}

		// create connection
		wac, err := whats.NewConn()
		if err != nil {
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write(badRequest)
			return
		}
		// check contact
		wphone := num
		var ci ContactInfo

		dd, er := wac.Exist(wphone + "@s.whatsapp.net")
		if er != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequest)
			return
		}
		ci.Exists = <-dd

		sd, fg := wac.GetStatus(wphone + "@s.whatsapp.net")
		if fg != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequest)
			return
		}
		ci.Status = <-sd

		a, b := wac.SubscribePresence(wphone + "@s.whatsapp.net")
		if b != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequest)
			return
		}
		ci.Online = <-a

		t, f := wac.GetProfilePicThumb(wphone + "@s.whatsapp.net")
		if f != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequest)
			return
		}
		ci.Thumb = <-t

		err = whats.SendMessageA(wac, device[0], num[0], msg[0])
		if err == whats.ErrConnecting {
			w.WriteHeader(http.StatusForbidden)
			w.Write(connFailed)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write(reqTimeout)
			return
		}

		w.WriteHeader(http.StatusOK)

		js, err := json.Marshal(ci)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
}

//MakeMessageHandlers make url handlers
func MakeContactHandler(r *mux.Router, n negroni.Negroni) {
	r.Handle("/check", n.With(negroni.Wrap(checkContact()))).Methods("GET").Name("checkContact")
}
