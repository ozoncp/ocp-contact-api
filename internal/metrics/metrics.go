package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	createCounter prometheus.Counter
	updateCounter prometheus.Counter
	removeCounter prometheus.Counter
)

func RegisterMetrics() {
	createCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ocp_contact_api_create_count_total",
		Help: "The total count of created contacts",
	})

	updateCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ocp_contact_api_update_count_total",
		Help: "The total count of updated contacts",
	})

	removeCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ocp_contact_api_remove_count_total",
		Help: "The total count of removed contacts",
	})
}

func CreateCounterInc() {
	createCounter.Inc()
}

func UpdateCounterInc() {
	updateCounter.Inc()
}

func RemoveCounterInc() {
	removeCounter.Inc()
}