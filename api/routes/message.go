package routes

import (
	"fmt"
	"net/http"

	"github.com/workstash/whapi/infrastructure/whats"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func sendMessage() http.Handler {
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
		msg, ok := r.URL.Query()["msg"]
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
		// send message

		num[0], err = whats.ValidateNum(num[0], wac, device[0])
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write(invalidNumber)
			return
		}

		fmt.Println("Enviando mensagem para o número ", num[0])

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
	})
}

//MakeMessageHandlers make url handlers
func MakeMessageHandlers(r *mux.Router, n negroni.Negroni) {
	r.Handle("/send", n.With(negroni.Wrap(sendMessage()))).Methods("GET").Name("sendMessage")
}
