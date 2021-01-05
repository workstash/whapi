package middleware

import (
	"net/http"
	"strconv"

	"github.com/workstash/whapi/helper/metrics"

	"github.com/urfave/negroni"
)

//Metrics to prometheus
func Metrics(mService metrics.Service) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		appMetric := metrics.NewHTTP(r.URL.Path, r.Method)
		appMetric.Started()
		next(w, r)
		res := w.(negroni.ResponseWriter)
		appMetric.Finished()
		appMetric.StatusCode = strconv.Itoa(res.Status())
		mService.SaveHTTP(appMetric)
	}
}
