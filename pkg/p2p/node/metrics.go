package node

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	IncomingConnectionCount prometheus.Counter
	OutgoingConnectionCount prometheus.Counter
	ConnectedPeerCount      prometheus.Counter
	DisconnectedPeerCount   prometheus.Counter
	BlockedPeerCount        prometheus.Counter
	CurrentBlockedPeerCount prometheus.Gauge
	BlockedPeers            *prometheus.GaugeVec
}

func newMetrics(registry prometheus.Registerer, namespace string) *metrics {
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
			Help:      "Number of connected peers.",
		}),
		DisconnectedPeerCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "disconnected_peer_count",
			Help:      "Number of disconnected peers.",
		}),
		BlockedPeerCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "blocked_peer_count",
			Help:      "Number of blocked peers.",
		}),
		CurrentBlockedPeerCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "current_blocked_peer_count",
			Help:      "Number of current blocked peers.",
		}),
		BlockedPeers: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "blocked_peers",
				Help:      "Blocked peers with reason.",
			},
			[]string{"peer_id", "reason"},
		),
	}

	registry.MustRegister(
		m.IncomingConnectionCount,
		m.OutgoingConnectionCount,
		m.ConnectedPeerCount,
		m.BlockedPeerCount,
		m.CurrentBlockedPeerCount,
		m.BlockedPeers,
	)

	return m
}
