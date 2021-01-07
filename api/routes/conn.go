package routes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/workstash/whapi/config"
	"github.com/workstash/whapi/infrastructure/qrcode"
	"github.com/workstash/whapi/infrastructure/whats"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func qrCode() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config := config.Main.API
		// validate data
		device, ok := r.URL.Query()["device"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequest)
			return
		}
		// create connection
		wac, err := whats.NewConn()
		if err != nil {
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write(reqTimeout)
			return
		}
		// response struct
		res := struct {
			QrCode string `json:"qrcode"`
		}{}

		// login
		qrCode, err := whats.Login(wac, device[0])
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write(connFailed)
			return
		}
		if qrCode != "" {
			if config.GenerateQrCode {
				qrCodeBytes, err := qrcode.GenerateQrCode(qrCode, config.QrCodeQuality, config.QrCodeSize)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(internalServerErr)
					return
				}
				qrCode = string(qrCodeBytes)
			}
			if config.EncodeBase64 {
				qrCode = base64.StdEncoding.EncodeToString([]byte(qrCode))
			}

			res.QrCode = qrCode
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(res); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(internalServerErr)
				return
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func close() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validate data
		device, ok := r.URL.Query()["device"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequest)
			return
		}
		// remove session
		err := os.Remove(fmt.Sprintf("%s/%s.gob", config.Main.API.SessionPath, device[0]))
		if err != nil {
			w.Write([]byte("session don't exists"))
		}

		w.WriteHeader(http.StatusOK)
	})
}

//MakeConnectionHandlers make url handlers
func MakeConnectionHandlers(r *mux.Router, n negroni.Negroni) {
	r.Handle("/qrcode", n.With(negroni.Wrap(qrCode()))).Methods("GET").Name("qrCode")
	r.Handle("/close", n.With(negroni.Wrap(close()))).Methods("GET").Name("closeconn")
}
