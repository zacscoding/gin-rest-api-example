package metric

import (
	"gin-rest-api-example/internal/config"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsProvider struct {
	Namespace string
	Subsystem string

	apiMetricsProvider apiMetricsProvider
}

type apiMetricsProvider struct {
	requestCounter *prometheus.CounterVec
	requestLatency *prometheus.SummaryVec
}

// RecordApiCount increases count of api request with given code, method, path labels
func (mp *MetricsProvider) RecordApiCount(code int, method, path string) {
	mp.apiMetricsProvider.requestCounter.WithLabelValues(strconv.Itoa(code), method, path).Inc()
}

// RecordApiLatency observes given elapsed mills with given code, method, path labels
func (mp *MetricsProvider) RecordApiLatency(code int, method, path string, elapsed time.Duration) {
	mills := float64(elapsed.Milliseconds())
	mp.apiMetricsProvider.requestLatency.WithLabelValues(strconv.Itoa(code), method, path).Observe(mills)
}

// NewMetricsProvider creates a new metrics provider to record metrics
func NewMetricsProvider(cfg *config.Config) *MetricsProvider {
	var (
		ns = cfg.MetricsConfig.Namespace
		ss = cfg.MetricsConfig.Subsystem
	)
	mp := MetricsProvider{
		Namespace: ns,
		Subsystem: ss,
		apiMetricsProvider: apiMetricsProvider{
			requestCounter: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: ns,
					Subsystem: ss,
					Name:      "api_request_count",
					Help:      "Total count of request",
				},
				[]string{"code", "method", "path"},
			),
			requestLatency: promauto.NewSummaryVec(
				prometheus.SummaryOpts{
					Namespace: ns,
					Subsystem: ss,
					Name:      "api_request_latency",
					Help:      "Elapsed time of request",
				},
				[]string{"code", "method", "path"},
			),
		},
	}
	return &mp
}
