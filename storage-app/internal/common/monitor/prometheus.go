package monitor

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total of http request",
		},
		[]string{"method", "endpoint", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_requests_duration",
			Help: "Request duration of http request",
		},
		[]string{"method", "endpoint", "status"},
	)
)

type PrometheusService interface {
	IncRequestTotal(_ context.Context, method, path string, status int)
	RecordRequestDuration(_ context.Context, method, path string, status int, duration float64)
}

type PrometheusServiceImpl struct {
	httpRequestTotal    *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
}

func NewPrometheusService(
	httpRequestTotal *prometheus.CounterVec,
	httpRequestDuration *prometheus.HistogramVec,
) PrometheusService {
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpRequestTotal)
	return &PrometheusServiceImpl{
		httpRequestTotal:    httpRequestTotal,
		httpRequestDuration: httpRequestDuration,
	}
}
func (p *PrometheusServiceImpl) IncRequestTotal(_ context.Context, method, path string, status int) {
	p.httpRequestTotal.WithLabelValues(method, path, http.StatusText(status)).Inc()
}

func (p *PrometheusServiceImpl) RecordRequestDuration(_ context.Context, method, path string, status int, duration float64) {
	p.httpRequestDuration.WithLabelValues(method, path, http.StatusText(status)).Observe(duration)
}
