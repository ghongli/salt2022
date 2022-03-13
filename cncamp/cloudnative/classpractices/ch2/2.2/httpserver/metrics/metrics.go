package metrics

import (
	"time"
	
	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricsNamespace = "httpserver"
)
var (
	funcLatency = CreateExecutionTimeMetric(MetricsNamespace, "Time spent.")
)

type (
	// ExecutionTimer measures execution time of a computation, split into major steps
	// usual usage pattern is: timer := NewExecutionTimer(...) ; compute ; timer.ObserveStep() ; ... ; timer.ObserveTotal()
	ExecutionTimer struct {
		histo *prometheus.HistogramVec
		start time.Time
		last time.Time
	}
)

// ObserveTotal measures the execution time from the creation of the ExecutionTimer.
func (t *ExecutionTimer) ObserveTotal() {
	(*t.histo).WithLabelValues("total").Observe(time.Since(t.start).Seconds())
}

// NewExecutionTimer provides a timer for admission latency; call ObserveXXX() on it to measure.
func NewExecutionTimer(histo *prometheus.HistogramVec) *ExecutionTimer {
	now := time.Now()
	return &ExecutionTimer{
		histo: histo,
		start: now,
		last:  now,
	}
}

// CreateExecutionTimeMetric prepares a new histogram labeled with execution step.
func CreateExecutionTimeMetric(ns, usage string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: ns,
		Name:      "execution_latency_seconds",
		Help:      usage,
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
	}, []string{"step"})
}

// NewTimer provides a timer for Updater's RunOnce execution.
func NewTimer() *ExecutionTimer {
	return NewExecutionTimer(funcLatency)
}

func Register() error {
	return prometheus.Register(funcLatency)
}