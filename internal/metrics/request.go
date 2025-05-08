package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// TotalRequestsCounter is a counter metric to track the total number of requests.
	// It increments every time a request is received. It is labeled with "path" (the target server path) and "method".
	TotalRequestsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "load_generation_system_total_requests_count", // Metric name
		},
		[]string{"path", "method"}, // Labels
	)

	// ProcessedRequestsCounter is a counter metric to track the number of processed requests.
	// It increments every time a request has been processed and will be labeled with "path" (the target server path),
	// "method", and "status" (the status code of the response).
	ProcessedRequestsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "load_generation_system_processed_requests_count", // Metric name
		},
		[]string{"path", "method", "status"}, // Labels
	)

	// RequestDurationSecondsHist is a histogram metric that tracks the duration of requests in seconds.
	// The histogram is labeled with "path" (the target server path), "method", and "status" (the status code).
	// It provides insight into how long requests take to complete, with predefined bucket ranges for different duration intervals.
	RequestDurationSecondsHist = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "load_generation_system_request_duration_seconds", // Metric name
			Buckets: []float64{ // Predefined bucket ranges for request durations in seconds.
				0.005, 0.01, 0.025, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9,
				1.0, 1.1, 1.2, 1.3, 1.4, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0, 10.0,
			},
		},
		[]string{"path", "method", "status"}, // Labels
	)
)
