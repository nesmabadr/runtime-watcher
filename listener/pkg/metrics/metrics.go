package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	listenerRequestDuration            = "watcher_listener_request_duration"
	listenerRequests                   = "watcher_listener_requests_total"
	listenerRequestErrors              = "watcher_listener_request_errors_total"
	listenerInflightRequests           = "watcher_listener_inflight_requests"
	listenerExceedingSizeLimitRequests = "watcher_listener_exceeding_size_limit_requests_total"
	listenerFailedVerificationRequests = "watcher_listener_failed_verification_requests_total"
	requestURILabel                    = "request_uri_label"
	listenerService                    = "listener"
	serverNameLabel                    = "server_name"
)

var (
	httpRequestDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{ //nolint:gochecknoglobals
		Name: listenerRequestDuration,
		Help: "Indicates the latency of each request in seconds",
	}, []string{serverNameLabel})
	httpRequestsCounter = prometheus.NewCounterVec(prometheus.CounterOpts{ //nolint:gochecknoglobals
		Name: listenerRequests,
		Help: "Indicates the number of requests",
	}, []string{serverNameLabel})
	httpRequestErrorsCounter = prometheus.NewCounterVec(prometheus.CounterOpts{ //nolint:gochecknoglobals
		Name: listenerRequestErrors,
		Help: "Indicates the number of failed requests",
	}, []string{serverNameLabel})
	HttpInflightRequestsGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{ //nolint:gochecknoglobals
		Name: listenerInflightRequests,
		Help: "Indicates the number of inflight requests",
	}, []string{serverNameLabel})
	httpRequestsExceedingSizeLimitCounter = prometheus.NewCounterVec(prometheus.CounterOpts{ //nolint:gochecknoglobals
		Name: listenerExceedingSizeLimitRequests,
		Help: "Indicates the number of requests exceeding size limit",
	}, []string{serverNameLabel})
	httpFailedVerificationRequests = prometheus.NewCounterVec(prometheus.CounterOpts{ //nolint:gochecknoglobals
		Name: listenerFailedVerificationRequests,
		Help: "Indicates the number of requests that failed verification",
	}, []string{serverNameLabel, requestURILabel})
)

func InitMetrics(metricsRegistry metrics.RegistererGatherer) {
	metricsRegistry.MustRegister(httpRequestDurationHistogram)
	metricsRegistry.MustRegister(httpRequestsCounter)
	metricsRegistry.MustRegister(httpRequestErrorsCounter)
	metricsRegistry.MustRegister(HttpInflightRequestsGauge)
	metricsRegistry.MustRegister(httpRequestsExceedingSizeLimitCounter)
	metricsRegistry.MustRegister(httpFailedVerificationRequests)
}

func UpdateMetrics(duration time.Duration) {
	recordHTTPRequestDuration(duration)
	recordHTTPRequests()
}

func recordHTTPRequestDuration(duration time.Duration) {
	httpRequestDurationHistogram.WithLabelValues(listenerService).Observe(duration.Seconds())
}

func recordHTTPRequests() {
	httpRequestsCounter.WithLabelValues(listenerService).Inc()
}

func RecordHTTPRequestErrors() {
	httpRequestErrorsCounter.WithLabelValues(listenerService).Inc()
}

func RecordHTTPInflightRequests(increaseBy float64) {
	HttpInflightRequestsGauge.WithLabelValues(listenerService).Add(increaseBy)
}

func RecordHTTPRequestExceedingSizeLimit() {
	httpRequestsExceedingSizeLimitCounter.WithLabelValues(listenerService).Inc()
}

func RecordHTTPFailedVerificationRequests(requestURI string) {
	httpFailedVerificationRequests.WithLabelValues(listenerService, requestURI).Inc()
}
