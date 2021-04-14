package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/workstash/whapi/api/middleware"
	"github.com/workstash/whapi/api/routes"
	"github.com/workstash/whapi/config"
	"github.com/workstash/whapi/helper/logger"
	"github.com/workstash/whapi/helper/metrics"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
)

func main() {
	defer logger.Close()
	metricService, err := metrics.NewPrometheusService()
	if err != nil {
		logger.Fatal(err.Error())
	}
	r := mux.NewRouter()
	n := negroni.New(
		negroni.HandlerFunc(middleware.Cors),
		negroni.HandlerFunc(middleware.Metrics(metricService)),
		negroni.NewLogger(),
	)

	//ROUTES
	//message
	routes.MakeConnectionHandlers(r, *n)
	routes.MakeMessageHandlers(r, *n)
	routes.MakeContactHandler(r, *n)

	http.Handle("/", r)
	http.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	f, err := os.OpenFile(fmt.Sprintf("%s/%s.log", config.Main.LoggerFile, time.Now().Format("2006-01-02")),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Println(err)
	}
	defer f.Close()
	negLog := log.New(f, "negLog: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + config.Main.API.Port,
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     negLog,
	}

	logger.Println("Server is running...")
	log.Printf("Port: %s", config.Main.API.Port)
	log.Println("Server is running...")
	logger.Fatal(srv.ListenAndServe().Error())
}
