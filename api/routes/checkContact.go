package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/workstash/whapi/config"
	"github.com/workstash/whapi/infrastructure/whats"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Structure to return JSON with Contacts Info

type payload struct {
	Status string `json:"status"`
	JID    string `json:"jid"`
}

type ContactInfo struct {
	Exists *payload `json:"exists"`
	Status *payload `json:"status"`
	Online *payload `json:"online"`
	Thumb  *payload `json:"thumb"`
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

		sessionPath := fmt.Sprintf("%s/%s.gob", config.Main.API.SessionPath, device[0])
		if err := whats.Auth(wac, sessionPath); err == nil {

			// check contact
			wphone := fmt.Sprintf("%v@s.whatsapp.net", num[0])

			//fmt.Println("num:", num)
			//fmt.Println("wphone:", wphone)

			var ci ContactInfo
			var pdd, psd, pa, pt *payload

			//------------- Exist ---------------
			dd, er := wac.Exist(wphone)
			if er != nil {
				fmt.Println("Failed Exist: ", er.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write(badRequest)
				return
			}
			if err := json.Unmarshal([]byte(<-dd), &pdd); err != nil {
				fmt.Println("Error unmarshalling Exist: ", err.Error())
				return
			}
			ci.Exists = pdd

			//------------- GetStatus ---------------
			sd, fg := wac.GetStatus(wphone)
			if fg != nil {
				fmt.Println("Failed GetStatus: ", fg.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write(badRequest)
				return
			}
			if err := json.Unmarshal([]byte(<-sd), &psd); err != nil {
				fmt.Println("Error unmarshalling GetStatus: ", err.Error())
				return
			}
			ci.Status = psd

			//------------- SubscribePresence ---------------
			a, b := wac.SubscribePresence(wphone)
			if b != nil {
				fmt.Println("Failed SubscribePresence: ", b.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write(badRequest)
				return
			}
			if err := json.Unmarshal([]byte(<-a), &pa); err != nil {
				fmt.Println("Error unmarshalling SubscribePresence: ", err.Error())
				return
			}

			ci.Online = pa

			//------------- GetProfilePicThumb ---------------
			t, f := wac.GetProfilePicThumb(wphone)
			if f != nil {
				fmt.Println("Failed GetProfilePicThumb: ", f.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write(badRequest)
				return
			}
			if err := json.Unmarshal([]byte(<-t), &pt); err != nil {
				fmt.Println("Error unmarshalling GetProfilePicThumb: ", err.Error())
				return
			}

			ci.Thumb = pt

			js, err := json.Marshal(ci)
			if err != nil {
				fmt.Println("Error marshalling Contact Info: ", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write(badRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		} else {
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write(reqTimeout)
			return
		}
	})
}

//MakeMessageHandlers make url handlers
func MakeContactHandler(r *mux.Router, n negroni.Negroni) {
	r.Handle("/check", n.With(negroni.Wrap(checkContact()))).Methods("GET").Name("checkContact")
}
