package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ActiveUsersGauge is a gauge metric that tracks the current number of active users.
	// Gauges represent a single numerical value that can go up or down. In this case, it is used to monitor
	// the number of users currently active in the system, which may increase or decrease over time based on user activity.
	ActiveUsersGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "load_generation_system_active_users", // Metric name
	})
)
