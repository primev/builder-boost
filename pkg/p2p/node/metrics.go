package node

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	IncomingConnectionCount prometheus.Counter
	OutgoingConnectionCount prometheus.Counter
	ConnectedPeerCount      prometheus.Counter
	DisconnectedPeerCount   prometheus.Counter
}

func newMetrics(registry prometheus.Registerer) *metrics {
	namespace := "builder_boost"
	subsystem := "p2p_node"

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
		ConnectedPeerCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "connected_peer_count",
			Help:      "Number of connected peer.",
		}),
		DisconnectedPeerCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "disconnected_peer_count",
			Help:      "Number of disconnected peer.",
		}),
	}

	registry.MustRegister(
		m.IncomingConnectionCount,
		m.OutgoingConnectionCount,
		m.ConnectedPeerCount,
		m.DisconnectedPeerCount,
	)

	return m
}
