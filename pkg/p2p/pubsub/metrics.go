package pubsub

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	PublishedMsgCount prometheus.Counter
	StreamedMsgCount  prometheus.Counter
	GossipedMsgCount  prometheus.Counter
}

func newMetrics(registry prometheus.Registerer) *metrics {
	namespace := "builder_boost"
	subsystem := "p2p_pubsub"

	m := &metrics{
		// outgoing metrics
		PublishedMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "published_msg_count",
			Help:      "Number of published message count.",
		}),
		StreamedMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "streamed_msg_count",
			Help:      "Number of streamed message count.",
		}),
		GossipedMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "gossiped_msg_count",
			Help:      "Number of gossiped message count.",
		}),
		// incoming metrics
	}

	registry.MustRegister(
		m.PublishedMsgCount,
		m.StreamedMsgCount,
		m.GossipedMsgCount,
	)

	return m
}
