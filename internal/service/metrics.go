package service

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	duration     *prometheus.HistogramVec
	addedUsers   prometheus.Counter
	deletedUsers prometheus.Counter
}

func newMetrics() *metrics {
	return &metrics{
		duration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "name_enricher_service",
				Subsystem: "",
				Name:      "db_resp_duration",
				Help:      "database response duration",
				Buckets:   []float64{0.0001, 0.0005, 0.001, 0.003, 0.005, 0.01, 0.05, 0.1, 1},
			}, []string{"operation_type"}),
		addedUsers: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: "name_enricher_service",
				Subsystem: "",
				Name:      "users_added_total",
				Help:      "total quantity of users that were added",
			}),
		deletedUsers: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: "name_enricher_service",
				Subsystem: "",
				Name:      "users_deleted_total",
				Help:      "total quantity of users that were deleted",
			}),
	}
}
