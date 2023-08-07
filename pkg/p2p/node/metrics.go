package node

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	IncomingConnectionCount prometheus.Counter
	OutgoingConnectionCount prometheus.Counter
}

func newMetrics(registry prometheus.Registerer) *metrics {
	namespace := "builder_boost"
	subsystem := "p2p"

	m := &metrics{
		IncomingConnectionCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "incoming_connection_count",
			Help:      "Number of initiated incoming libp2p connections.",
		}),
		OutgoingConnectionCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "outgoing_connection_count",
			Help:      "Number of handled outgoing libp2p connections.",
		}),
	}

	registry.MustRegister(
		m.IncomingConnectionCount,
		m.OutgoingConnectionCount,
	)

	return m
}
