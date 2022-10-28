package metrics

import (
	"github.com/penglongli/gin-metrics/ginmetrics"
)

// var totalRequests = prometheus.NewCounterVec(
// 	prometheus.CounterOpts{
// 		Name: "requests_total",
// 		Help: "Number of requests.",
// 	},
// 	[]string{"method", "path"},
// )

// var responseStatus = prometheus.NewCounterVec(
// 	prometheus.CounterOpts{
// 		Name: "response_status",
// 		Help: "Status of HTTP response",
// 	},
// 	[]string{"status"},
// )

// var readCache = prometheus.NewCounter(
// 	prometheus.CounterOpts{
// 		Name: "use_cache",
// 		Help: "Number of cache extractions",
// 	},
// )

type metrics struct {
	monitor *ginmetrics.Monitor
}

func NewMonitor() *ginmetrics.Monitor {
	monitor := ginmetrics.GetMonitor()
	monitor.SetMetricPath("/metrics")

	metrics := new(metrics)
	metrics.monitor = monitor
	metrics.makeMetrics()

	return metrics.monitor
}

func (m *metrics) makeMetrics() {
	totalRequests := ginmetrics.Metric{
		Type:        ginmetrics.Counter,
		Name:        "requests_total",
		Description: "Number of requests.",
		Labels:      []string{"method", "path"},
	}
	_ = m.monitor.AddMetric(&totalRequests)

	responseStatus := ginmetrics.Metric{
		Type:        ginmetrics.Counter,
		Name:        "response_status",
		Description: "Status of HTTP response.",
		Labels:      []string{"status"},
	}
	_ = m.monitor.AddMetric(&responseStatus)

	readCache := ginmetrics.Metric{
		Type:        ginmetrics.Counter,
		Name:        "read_cache",
		Description: "Number of cache extractions",
		Labels:      nil,
	}
	_ = m.monitor.AddMetric(&readCache)
}
