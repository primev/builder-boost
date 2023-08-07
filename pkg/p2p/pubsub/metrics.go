package pubsub

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	// outgoing metrics
	PublishedMsgCount prometheus.Counter
	StreamedMsgCount  prometheus.Counter
	GossipedMsgCount  prometheus.Counter
	// incoming metrics
	ApproveMsgCount     prometheus.Counter
	PingMsgCount        prometheus.Counter
	PongMsgCount        prometheus.Counter
	GetVersionMsgCount  prometheus.Counter
	VersionMsgCount     prometheus.Counter
	GetPeerListMsgCount prometheus.Counter
	PeerListMsgCount    prometheus.Counter
	SignatureMsgCount   prometheus.Counter
	BlockKeyMsgCount    prometheus.Counter
	BundleMsgCount      prometheus.Counter
	PreconfBidMsgCount  prometheus.Counter
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
			Help:      "Number of published messages count.",
		}),
		StreamedMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "streamed_msg_count",
			Help:      "Number of streamed messages count.",
		}),
		GossipedMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "gossiped_msg_count",
			Help:      "Number of gossiped messages count.",
		}),
		// incoming metrics
		ApproveMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "approve_msg_count",
			Help:      "Number of incoming approve messages count.",
		}),
		PingMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "ping_msg_count",
			Help:      "Number of incoming ping messages count.",
		}),
		PongMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "pong_msg_count",
			Help:      "Number of incoming pong messages count.",
		}),
		GetVersionMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "getversion_msg_count",
			Help:      "Number of incoming getversion messages count.",
		}),
		VersionMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "version_msg_count",
			Help:      "Number of incoming version messages count.",
		}),
		GetPeerListMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "getpeerlist_msg_count",
			Help:      "Number of incoming getpeerlist messages count.",
		}),
		PeerListMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "peerlist_msg_count",
			Help:      "Number of incoming peerlist messages count.",
		}),
		SignatureMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "signature_msg_count",
			Help:      "Number of incoming signature messages count.",
		}),
		BlockKeyMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "blockkey_msg_count",
			Help:      "Number of incoming blockkey messages count.",
		}),
		BundleMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "bundle_msg_count",
			Help:      "Number of incoming bundle messages count.",
		}),
		PreconfBidMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "preconfbid_msg_count",
			Help:      "Number of incoming preconfbid messages count.",
		}),
	}

	registry.MustRegister(
		// out
		m.PublishedMsgCount,
		m.StreamedMsgCount,
		m.GossipedMsgCount,
		// in
		m.ApproveMsgCount,
		m.PingMsgCount,
		m.PongMsgCount,
		m.GetVersionMsgCount,
		m.VersionMsgCount,
		m.GetPeerListMsgCount,
		m.PeerListMsgCount,
		m.SignatureMsgCount,
		m.BlockKeyMsgCount,
		m.BundleMsgCount,
		m.PreconfBidMsgCount,
	)

	return m
}
