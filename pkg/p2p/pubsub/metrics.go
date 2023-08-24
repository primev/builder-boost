package pubsub

import (
	"github.com/primev/builder-boost/pkg/p2p/commons"
	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	// the current number of connected peers in the p2p network
	ApprovedPeerCount prometheus.Gauge
	BuilderPeerCount  prometheus.Gauge
	SearcherPeerCount prometheus.Gauge
	// outgoing metrics b2b
	PublishedB2BMsgCount prometheus.Counter
	StreamedB2BMsgCount  prometheus.Counter
	GossipedB2BMsgCount  prometheus.Counter
	// outgoing metrics b2s
	PublishedB2SMsgCount prometheus.Counter
	StreamedB2SMsgCount  prometheus.Counter
	GossipedB2SMsgCount  prometheus.Counter
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
	BidMsgCount         prometheus.Counter
	CommitmentMsgCount  prometheus.Counter
	// rtt measurements are used to determine the latency in communication between peers
	LatencyPeers *prometheus.GaugeVec
	// the score values assigned by this peer address to other peers
	ScorePeers *prometheus.GaugeVec
	// the times when the peers joined to the network
	JoinDatePeers *prometheus.GaugeVec
}

func newMetrics(registry prometheus.Registerer, namespace string, peerType commons.PeerType) *metrics {
	subsystem := "p2p_pubsub"

	m := &metrics{
		// system metrics
		ApprovedPeerCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "approved_peer_count",
			Help:      "Number of approved peers.",
		}),

		BuilderPeerCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "approved_builder_peer_count",
			Help:      "Number of approved peers.",
		}),

		SearcherPeerCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "approved_searcher_peer_count",
			Help:      "Number of approved peers.",
		}),

		// outgoing metrics b2b
		PublishedB2BMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "published_b2b_msg_count",
			Help:      "Number of published b2b messages count.",
		}),
		StreamedB2BMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "streamed_b2b_msg_count",
			Help:      "Number of streamed b2b messages count.",
		}),
		GossipedB2BMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "gossiped_b2b_msg_count",
			Help:      "Number of gossiped b2b messages count.",
		}),

		// outgoing metrics b2s
		PublishedB2SMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "published_b2s_msg_count",
			Help:      "Number of published b2s messages count.",
		}),
		StreamedB2SMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "streamed_b2s_msg_count",
			Help:      "Number of streamed b2s messages count.",
		}),
		GossipedB2SMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "gossiped_b2s_msg_count",
			Help:      "Number of gossiped b2s messages count.",
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
		BidMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "bid_msg_count",
			Help:      "Number of incoming bid messages count.",
		}),
		CommitmentMsgCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "commitment_msg_count",
			Help:      "Number of incoming commitment messages count.",
		}),

		// rtt metrics
		LatencyPeers: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "peer_latency_milliseconds",
			Help:      "Peer to peer latency.",
		},
			[]string{"peer_id"},
		),

		// score metrics
		ScorePeers: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "peer_score",
			Help:      "Peer to peer score.",
		},
			[]string{"peer_id"},
		),

		// join date metrics
		JoinDatePeers: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "peer_join_date",
			Help:      "Unix time peer join date.",
		},
			[]string{"peer_id"},
		),
	}

	switch peerType {
	case commons.Builder:
		registry.MustRegister(
			// system
			m.ApprovedPeerCount,
			m.BuilderPeerCount,
			m.SearcherPeerCount,
			// out b2b
			m.PublishedB2BMsgCount,
			m.StreamedB2BMsgCount,
			m.GossipedB2BMsgCount,
			// out b2s
			m.PublishedB2SMsgCount,
			m.StreamedB2SMsgCount,
			m.GossipedB2SMsgCount,
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
			m.BidMsgCount,
			m.CommitmentMsgCount,
			// rtt
			m.LatencyPeers,
			// score
			m.ScorePeers,
			// join date
			m.JoinDatePeers,
		)

	case commons.Searcher:
		registry.MustRegister(
			// system
			m.ApprovedPeerCount,
			// out b2s
			m.PublishedB2SMsgCount,
			m.StreamedB2SMsgCount,
			m.GossipedB2SMsgCount,
			// in
			m.ApproveMsgCount,
			m.PingMsgCount,
			m.PongMsgCount,
			m.GetVersionMsgCount,
			m.VersionMsgCount,
			m.GetPeerListMsgCount,
			m.PeerListMsgCount,
			//TODO I'm not sure if Searchers will be take bid from the builders if not remove
			m.BidMsgCount,
			m.CommitmentMsgCount,
			// rtt
			m.LatencyPeers,
			// score
			m.ScorePeers,
			// join date
			m.JoinDatePeers,
		)
	}

	return m
}
