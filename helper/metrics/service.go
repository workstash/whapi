package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

//Prometheus implements Service interface
type Prometheus struct {
	httpRequestHistogram *prometheus.HistogramVec
}

//NewPrometheusService create a new prometheus service
func NewPrometheusService() (*Prometheus, error) {
	http := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler", "method", "code"})

	s := &Prometheus{
		httpRequestHistogram: http,
	}
	err := prometheus.Register(s.httpRequestHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	return s, nil
}

//SaveHTTP send metrics to server
func (s *Prometheus) SaveHTTP(h *HTTP) {
	s.httpRequestHistogram.WithLabelValues(h.Handler, h.Method, h.StatusCode).Observe(h.Duration)
}
