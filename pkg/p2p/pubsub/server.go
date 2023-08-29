package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/lthibault/log"
	"github.com/primev/builder-boost/pkg/p2p/commons"
	"github.com/primev/builder-boost/pkg/p2p/commons/messages"
	"github.com/primev/builder-boost/pkg/p2p/commons/signer"
	"github.com/primev/builder-boost/pkg/p2p/config"
	"github.com/primev/builder-boost/pkg/p2p/message"
	"github.com/primev/builder-boost/pkg/p2p/stream"
	"github.com/primev/builder-boost/pkg/rollup"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	errInvalidBuilderOpts  = errors.New("invalid builder options!")
	errInvalidSearcherOpts = errors.New("invalid searcher options!")
	errInvalidPeerType     = errors.New("invalid peer type!")
)

var (
	streamTimeout  = time.Second * 10
	connectTimeout = time.Second * 10
)

const (
	BuilderTopicsCount       = 2
	SearcherTopicsCount      = 1
	BuilderComChannelsCount  = 4
	SearcherComChannelsCount = 2
)

type streamFunc func(peer.ID, message.OutboundMessage) error

type PubSubIO struct {
	BuilderServer
	SearcherServer
}

type BuilderServer struct {
	Server
	topic_b2b  *pubsub.Topic
	stream_b2b stream.Stream

	comChannels map[commons.ComChannels]chan messages.PeerMsg

	ready chan struct{}
}

type SearcherServer struct {
	Server
	topic_b2s  *pubsub.Topic
	stream_b2s stream.Stream

	comChannels map[commons.ComChannels]chan messages.PeerMsg

	ready chan struct{}
}

type Server struct {
	ctx         context.Context
	log         log.Logger
	selfType    commons.PeerType
	selfID      peer.ID
	cfg         *config.Config
	host        host.Host
	connTrackCh chan commons.ConnectionEvent
	token       []byte
	address     common.Address
	rollup      rollup.Rollup
	imb         message.InboundMsgBuilder
	omb         message.OutboundMsgBuilder
	apm         *approvedPeersMap
	metrics     *metrics
}

func New(
	ctx context.Context,
	cfg *config.Config,
	log log.Logger,
	peerType commons.PeerType,
	host host.Host,
	registry *prometheus.Registry,
	connTrackCh chan commons.ConnectionEvent,
	token []byte,
	address common.Address,
	rollup rollup.Rollup,
	topics map[commons.Topic]*pubsub.Topic,
	imb message.InboundMsgBuilder,
	omb message.OutboundMsgBuilder,
	comChannels map[commons.ComChannels]chan messages.PeerMsg,
) (*PubSubIO, error) {

	switch peerType {
	case commons.Builder:
		if (len(topics) != BuilderTopicsCount) || (len(comChannels) != (BuilderComChannelsCount + SearcherComChannelsCount)) {
			return nil, errInvalidBuilderOpts
		}
	case commons.Searcher:
		if (len(topics) != SearcherTopicsCount) || (len(comChannels) != SearcherComChannelsCount) {
			return nil, errInvalidSearcherOpts
		}
	default:
		return nil, errInvalidPeerType
	}

	server := new(Server)
	apm := newApprovedPeersMap()

	// register pubsub metrics
	metrics := newMetrics(registry, cfg.MetricsNamespace(), peerType)

	server = &Server{
		ctx:         ctx,
		cfg:         cfg,
		log:         log,
		selfID:      host.ID(),
		selfType:    peerType,
		host:        host,
		connTrackCh: connTrackCh,
		token:       token,
		address:     address,
		rollup:      rollup,
		imb:         imb,
		omb:         omb,
		apm:         apm,
		metrics:     metrics,
	}

	var (
		builderServer  *BuilderServer
		searcherServer *SearcherServer
	)

	// make sure there is no empty value
	builderServer = &BuilderServer{}
	searcherServer = &SearcherServer{}

	switch peerType {
	case commons.Builder:
		builderServer = &BuilderServer{
			Server:      *server,
			topic_b2b:   topics[commons.TopicB2B],
			comChannels: comChannels,
			ready:       make(chan struct{}),
		}

		builderServer.stream_b2b = stream.New(
			host,
			protocol.ID(cfg.PeerStreamB2B()),
			builderServer.peerStreamHandlerB2B,
		)

		searcherServer = &SearcherServer{
			Server:      *server,
			topic_b2s:   topics[commons.TopicB2S],
			comChannels: comChannels,
			ready:       make(chan struct{}),
		}

		searcherServer.stream_b2s = stream.New(
			host,
			protocol.ID(cfg.PeerStreamB2S()),
			searcherServer.peerStreamHandlerB2S,
		)
	case commons.Searcher:
		searcherServer = &SearcherServer{
			Server:      *server,
			topic_b2s:   topics[commons.TopicB2S],
			comChannels: comChannels,
			ready:       make(chan struct{}),
		}

		searcherServer.stream_b2s = stream.New(
			host,
			protocol.ID(cfg.PeerStreamB2S()),
			searcherServer.peerStreamHandlerB2S,
		)
	default:

	}

	switch peerType {
	case commons.Builder:
		// event tracking
		go server.events(connTrackCh, builderServer.Stream, searcherServer.Stream)
		// start RTT test | don't use libp2p ping
		go server.latencyUpdater(builderServer.Stream, searcherServer.Stream)
	case commons.Searcher:
		// event tracking
		go server.events(connTrackCh, searcherServer.Stream)
		// start RTT test | don't use libp2p ping
		go server.latencyUpdater(searcherServer.Stream)
	default:
	}

	// start score test
	go server.scoreUpdater()

	switch peerType {
	case commons.Builder:
		var onceBuilder sync.Once
		var onceSearcher sync.Once

		go builderServer.builderProtocol(onceBuilder)
		go searcherServer.searcherProtocol(onceSearcher)

		<-builderServer.ready
		<-searcherServer.ready

		close(builderServer.ready)
		close(searcherServer.ready)
	case commons.Searcher:
		var onceSearcher sync.Once

		go searcherServer.searcherProtocol(onceSearcher)

		<-searcherServer.ready

		close(searcherServer.ready)
	default:
	}

	return &PubSubIO{
		BuilderServer:  *builderServer,
		SearcherServer: *searcherServer,
	}, nil
}

// base builder pubsub procotol
func (bs *BuilderServer) builderProtocol(once sync.Once) {
	bs.log.With(log.F{
		"caller":  commons.GetCallerName(),
		"date":    commons.GetNow(),
		"service": "p2p pubsub",
	}).Info("starting builder pubsub protocol...")

	sub, err := bs.topic_b2b.Subscribe()
	if err != nil {
		bs.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p pubsub",
		}).Fatal(err)
	}

	defer sub.Cancel()
	for {
		select {
		case <-bs.ctx.Done():
			bs.log.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p pubsub",
			}).Warn("stopped builder pubsub protocol...")

			return
		default:
		}

		once.Do(func() {
			bs.ready <- struct{}{}
		})

		msg, err := sub.Next(bs.ctx)
		if err != nil {
			bs.log.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p pubsub protocol reader",
			}).Error(err)
			continue
		}

		// check self peer.ID
		if strings.Compare(bs.selfID.String(), msg.ReceivedFrom.String()) == 0 {
			continue
		}

		inMsg, err := bs.imb.Parse(msg.ReceivedFrom, msg.Data)
		if err != nil {
			bs.log.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p pubsub protocol inbound msg parse",
			}).Error(err)
			continue
		}

		// if not in approved peers try to authenticate
		if !bs.apm.InBuilderPeers(msg.ReceivedFrom) {
			go func() {
				bs.log.With(log.F{
					"caller":  commons.GetCallerName(),
					"date":    commons.GetNow(),
					"service": "p2p pubsub",
					"op":      inMsg.Op(),
					"peer":    inMsg.Peer(),
				}).Debug("unverified peer message")

				switch inMsg.Op() {
				case message.Approve:
					bs.optApprove(bs.Stream, inMsg.Peer(), inMsg.Bytes(), true)
				default:
					bs.log.With(log.F{
						"caller":  commons.GetCallerName(),
						"date":    commons.GetNow(),
						"service": "p2p pubsub builder",
						"op":      inMsg.Op(),
						"peer":    inMsg.Peer(),
					}).Warn("unknown approve option!")
				}
			}()
			continue
		}

		go func() {
			bs.log.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p pubsub",
				"op":      inMsg.Op(),
				"peer":    inMsg.Peer(),
			}).Debug("verified peer message")

			switch inMsg.Op() {
			// pass auth option in this side for now
			case message.Approve:
			// create pong message and publish to show you're alive
			case message.Ping:
				bs.optPing(bs.Stream, inMsg.Peer(), inMsg.Bytes())
			// it can be use validate peer is alive
			case message.Pong:
				bs.optPong(inMsg.Peer(), inMsg.Bytes())
			// send version
			case message.GetVersion:
				bs.optGetVersion(bs.Stream, inMsg.Peer())
			// store peers version
			case message.Version:
				bs.optVersion(inMsg.Peer(), inMsg.Bytes())
			// send approved peers
			case message.GetPeerList:
				bs.optGetPeerList(bs.Stream, inMsg.Peer())
			// it can be use for connect to peers
			case message.PeerList:
				bs.optPeerList(inMsg.Peer(), inMsg.Bytes())
			// no permission over publish
			case message.Signature:
			// no permission over publish
			case message.BlockKey:
			// handle the incoming encrypted transactions
			case message.Bundle:
				bs.optBundle(inMsg.Peer(), inMsg.Bytes())
			// handle the incoming preconf bids
			case message.PreconfBid:
				bs.optPreconfBid(inMsg.Peer(), inMsg.Bytes())
			default:
				bs.log.With(log.F{
					"caller":  commons.GetCallerName(),
					"date":    commons.GetNow(),
					"service": "p2p pubsub",
					"op":      inMsg.Op(),
					"peer":    inMsg.Peer(),
				}).Warn("unknown option!")
			}
		}()
	}
}

// base builder pussub procotol
func (ss *SearcherServer) searcherProtocol(once sync.Once) {
	ss.log.With(log.F{
		"caller":  commons.GetCallerName(),
		"date":    commons.GetNow(),
		"service": "p2p pussub",
	}).Info("starting searcher pussub protocol...")

	sub, err := ss.topic_b2s.Subscribe()
	if err != nil {
		ss.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p pubsub",
		}).Fatal(err)
	}

	defer sub.Cancel()
	for {
		select {
		case <-ss.ctx.Done():
			ss.log.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p pubsub",
			}).Warn("stopped searcher pussub protocol...")

			return
		default:
		}

		once.Do(func() {
			ss.ready <- struct{}{}
		})

		msg, err := sub.Next(ss.ctx)
		if err != nil {
			ss.log.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p pubsub protocol reader",
			}).Error(err)
			continue
		}

		// check self peer.ID
		if strings.Compare(ss.selfID.String(), msg.ReceivedFrom.String()) == 0 {
			continue
		}

		// builder to builder communication is not allowed
		if (ss.selfType == commons.Builder) && (ss.apm.InBuilderPeers(msg.ReceivedFrom)) {
			continue
		}

		inMsg, err := ss.imb.Parse(msg.ReceivedFrom, msg.Data)
		if err != nil {
			ss.log.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p pubsub protocol",
			}).Error(err)
			continue
		}

		// if not in approved peers try to approve
		if ((ss.selfType == commons.Builder) && !(ss.apm.InSearcherPeers(msg.ReceivedFrom))) ||
			((ss.selfType == commons.Searcher) && !(ss.apm.InBuilderPeers(msg.ReceivedFrom))) {
			//if !ss.apm.InBuilderPeers(msg.ReceivedFrom) {

			go func() {
				ss.log.With(log.F{
					"caller":  commons.GetCallerName(),
					"date":    commons.GetNow(),
					"service": "p2p pubsub",
					"op":      inMsg.Op(),
					"peer":    inMsg.Peer(),
				}).Debug("unverified peer message")

				switch inMsg.Op() {
				case message.Approve:
					ss.optApprove(ss.Stream, inMsg.Peer(), inMsg.Bytes(), true)
				default:
					ss.log.With(log.F{
						"caller":  commons.GetCallerName(),
						"date":    commons.GetNow(),
						"service": "p2p pubsub searcher",
						"op":      inMsg.Op(),
						"peer":    inMsg.Peer(),
					}).Warn("unknown approve option!")
				}
			}()
			continue
		}

		go func() {
			ss.log.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p pubsub",
				"op":      inMsg.Op(),
				"peer":    inMsg.Peer(),
			}).Debug("verified peer message")

			switch inMsg.Op() {
			// pass auth option in this side for now
			case message.Approve:
			// create pong message and publish to show you're alive
			case message.Ping:
				ss.optPing(ss.Stream, inMsg.Peer(), inMsg.Bytes())
			// it can be use validate peer is alive
			case message.Pong:
				ss.optPong(inMsg.Peer(), inMsg.Bytes())
			// send version
			case message.GetVersion:
				ss.optGetVersion(ss.Stream, inMsg.Peer())
			// store peers version
			case message.Version:
				ss.optVersion(inMsg.Peer(), inMsg.Bytes())
			// send approved peers
			case message.GetPeerList:
				ss.optGetPeerList(ss.Stream, inMsg.Peer())
			// it can be use for connect to peers
			case message.PeerList:
				ss.optPeerList(inMsg.Peer(), inMsg.Bytes())
			// searcher<>builder communication test
			case message.Bid:
				ss.optBid(inMsg.Peer(), inMsg.Bytes())
			// searcher<>builder communication test
			case message.Commitment:
				ss.optCommitment(inMsg.Peer(), inMsg.Bytes())

			default:
				ss.log.With(log.F{
					"caller":  commons.GetCallerName(),
					"date":    commons.GetNow(),
					"service": "p2p pubsub",
					"op":      inMsg.Op(),
					"peer":    inMsg.Peer(),
				}).Warn("unknown option!")
			}
		}()
	}
}

// publish message on topic
func (bs *BuilderServer) Publish(msg message.OutboundMessage) error {
	defer bs.metrics.PublishedB2BMsgCount.Inc()

	msgBytes, err := msg.MarshalJSON()
	if err != nil {
		return err
	}

	return bs.topic_b2b.Publish(bs.ctx, msgBytes)
}

// stream message for specific peer
func (bs *BuilderServer) Stream(peerID peer.ID, msg message.OutboundMessage) error {
	defer bs.metrics.StreamedB2BMsgCount.Inc()

	msgBytes, err := msg.MarshalJSON()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(bs.ctx, streamTimeout)
	defer cancel()

	return bs.stream_b2b.Send(ctx, peerID, msgBytes)
}

// stream message for gossip peers
func (bs *BuilderServer) Gossip(msg message.OutboundMessage) error {
	defer bs.metrics.GossipedB2BMsgCount.Inc()

	msgBytes, err := msg.MarshalJSON()
	if err != nil {
		return err
	}

	// TODO use optimization to fast data distribution
	peers := bs.apm.GetGossipPeers()

	for peerID, _ := range peers {
		ctx, cancel := context.WithTimeout(bs.ctx, streamTimeout)
		err = bs.stream_b2b.Send(ctx, peerID, msgBytes)
		if err != nil {
			bs.log.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p pubsub gossip",
				"peer id": peerID,
			}).Error(err)
		}
		cancel()
	}

	return nil
}

// publish message on topic
func (ss *SearcherServer) Publish(msg message.OutboundMessage) error {
	defer ss.metrics.PublishedB2SMsgCount.Inc()

	msgBytes, err := msg.MarshalJSON()
	if err != nil {
		return err
	}

	return ss.topic_b2s.Publish(ss.ctx, msgBytes)
}

// stream message for specific peer
func (ss *SearcherServer) Stream(peerID peer.ID, msg message.OutboundMessage) error {
	defer ss.metrics.StreamedB2SMsgCount.Inc()

	msgBytes, err := msg.MarshalJSON()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ss.ctx, streamTimeout)
	defer cancel()

	return ss.stream_b2s.Send(ctx, peerID, msgBytes)
}

// stream message for gossip peers
func (ss *SearcherServer) Gossip(msg message.OutboundMessage) error {
	defer ss.metrics.GossipedB2SMsgCount.Inc()

	msgBytes, err := msg.MarshalJSON()
	if err != nil {
		return err
	}

	// TODO use optimization to fast data distribution
	peers := ss.apm.GetGossipPeers()

	for peerID, _ := range peers {
		ctx, cancel := context.WithTimeout(ss.ctx, streamTimeout)
		err = ss.stream_b2s.Send(ctx, peerID, msgBytes)
		if err != nil {
			ss.log.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p pubsub gossip",
				"peer id": peerID,
			}).Error(err)
		}
		cancel()
	}

	return nil
}

func (server *Server) optApprove(stream streamFunc, cpeer peer.ID, bytes []byte, sendauth bool) {
	defer server.metrics.ApproveMsgCount.Inc()

	var am = new(messages.ApproveMsg)
	err := json.Unmarshal(bytes, &am)

	newSigner := signer.New()
	valid, address, err := newSigner.Verify(am.Sig, am.GetUnsignedMessage())
	if err != nil {
		server.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p pubsub approve",
		}).Error(err)

		// terminate the unexpected connection
		server.host.Network().ClosePeer(cpeer)
		return
	}

	// verify signature
	if !valid {
		server.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p pubsub approve",
		}).Error(errors.New("not valid signature"))

		// terminate the unexpected connection
		server.host.Network().ClosePeer(cpeer)
		return
	}

	// verify the correct builder eth address packaging
	if !commons.BytesCompare(address.Bytes(), am.Address.Bytes()) {
		server.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p pubsub approve",
		}).Error(errors.New("wrong adress packaging"))

		// terminate the unexpected connection
		server.host.Network().ClosePeer(cpeer)
		return
	}

	// terminate clone peer connections
	if commons.BytesCompare(server.address.Bytes(), am.Address.Bytes()) {
		server.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p pubsub approve",
		}).Error(errors.New("clone builder detected"))

		// terminate the unexpected connection
		server.host.Network().ClosePeer(cpeer)
		return
	}

	// verify the correct builder peer address packaging
	if strings.Compare(cpeer.String(), am.Peer.String()) != 0 {
		server.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p pubsub approve",
		}).Error(errors.New("wrong peer packaging"))

		// terminate the unexpected connection
		server.host.Network().ClosePeer(cpeer)
		return
	}

	var (
		stake *big.Int
	)

	switch am.PeerType {
	// If a peer wants to be approved in builder mode, we can use rollup to
	// determine whether it is a registered builder or not
	case commons.Builder:
		// get stake amount
		stake, err = server.rollup.GetMinimalStake(address)
		if err != nil {
			server.apm.DelPeer(cpeer)
			return
		}

		// check builder stake amount
		if stake.Cmp(server.cfg.MinimalStake()) < 0 {
			// terminate the unexpected connection
			server.host.Network().ClosePeer(cpeer)
			server.apm.DelPeer(cpeer)

			server.log.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p pubsub approve",
				"peer":    cpeer,
			}).Warn("not enough stake")
			return
		}

	// TODO: when rollup support arrives for searchers, the validation here will change
	// use a whitelist temporarily
	case commons.Searcher:
		var inTrustedSearchers = false
		for _, v := range config.SearcherPeerIDs {
			pid, _ := peer.Decode(v)
			if cpeer == pid {
				inTrustedSearchers = true
			}
		}

		if !inTrustedSearchers {
			server.apm.DelPeer(cpeer)
			return
		}

	default:
		// terminate the unexpected connection
		server.host.Network().ClosePeer(cpeer)

		server.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p pubsub approve",
			"peer":    cpeer,
		}).Warn("unexpected approve message!")
		return

	}

	for p, i := range server.apm.GetPeers() {
		switch server.host.Network().Connectedness(p) {
		case network.NotConnected:
			server.apm.DelPeer(p)
		case network.CannotConnect:
			server.apm.DelPeer(p)
		}

		// close the old connection if builder peer to establish a new one without closing it
		if commons.BytesCompare(i.address.Bytes(), am.Address.Bytes()) {
			server.apm.DelPeer(p)
			if server.host.Network().Connectedness(p) == network.Connected {
				server.host.Network().ClosePeer(p)
			}
		}
	}

	addrInfo := server.host.Peerstore().PeerInfo(cpeer)

	server.apm.AddPeer(cpeer)
	joinDate := time.Now()

	server.apm.SetPeerInfoPeerType(cpeer, am.PeerType)
	server.apm.SetPeerInfoJoinDate(cpeer, joinDate)
	server.apm.SetPeerInfoAddress(cpeer, address)
	server.apm.SetPeerInfoStake(cpeer, stake)
	server.apm.SetPeerInfoAddrs(cpeer, addrInfo.Addrs)
	server.apm.SetPeerInfoUUID(cpeer, uuid.New())
	// temporary gossip
	server.apm.SetPeerInfoGossip(cpeer, true)

	server.metrics.JoinDatePeers.WithLabelValues(cpeer.String()).Set(float64(joinDate.Unix()))

	if sendauth {
		// create approve message and stream this message for new peer
		msg, err := server.omb.Approve(server.token)
		if err != nil {
			return
		}

		stream(cpeer, msg)
	}
}

// Create a 'pong' message and stream it to the peer that received the 'ping' message.
func (server *Server) optPing(stream streamFunc, cpeer peer.ID, uuidBytes []byte) {
	defer server.metrics.PingMsgCount.Inc()

	msg, err := server.omb.Pong(uuidBytes)
	if err != nil {
		return
	}

	stream(cpeer, msg)
}

// Set the latency value by checking the UUID and peer match in the received 'pong' message.
func (server *Server) optPong(cpeer peer.ID, uuidBytes []byte) {
	defer server.metrics.PongMsgCount.Inc()

	var newUUID uuid.NullUUID
	err := newUUID.UnmarshalBinary(uuidBytes)
	if err != nil {
		return
	}

	info := server.apm.GetPeerInfo(cpeer)

	if info.uuid == newUUID.UUID {
		server.apm.SetPeerInfoPongTime(cpeer, time.Now().UnixNano())
		info := server.apm.GetPeerInfo(cpeer)

		// calculate latency value
		latency := time.Duration(info.pongTime - info.pingTime)

		// set latency value
		server.apm.SetPeerInfoLatency(cpeer, latency)

		// set metric values
		latencyMs := float64(latency.Microseconds()) / 1000.0
		server.metrics.LatencyPeers.WithLabelValues(cpeer.String()).Set(latencyMs)

		// change uuid for next ping-pong
		server.apm.SetPeerInfoUUID(cpeer, uuid.New())
	} else {
		// terminate the peer that engages in unexpected data exchange
		server.host.Network().ClosePeer(cpeer)
	}
}

// Create a 'version' message and stream it to the peer that received the 'getversion' message.
func (server *Server) optGetVersion(stream streamFunc, cpeer peer.ID) {
	defer server.metrics.GetVersionMsgCount.Inc()

	msg, err := server.omb.Version(server.cfg.Version())
	if err != nil {
		return
	}

	stream(cpeer, msg)
}

// store peer version info in approved peers map
func (server *Server) optVersion(cpeer peer.ID, bytes []byte) {
	defer server.metrics.VersionMsgCount.Inc()

	server.apm.SetPeerInfoVersion(cpeer, bytes)
}

// Create a 'peerlist' message and stream it to the peer that received the 'getpeerlist' message.
func (server *Server) optGetPeerList(stream streamFunc, cpeer peer.ID) {
	defer server.metrics.GetPeerListMsgCount.Inc()

	// whether the requester is a searcher or a builder, always send the builder list
	msg, err := server.omb.PeerList(server.apm.GetBuilderPeerAddrs())
	if err != nil {
		return
	}

	stream(cpeer, msg)
}

// get peerlist from other peers
func (server *Server) optPeerList(cpeer peer.ID, bytes []byte) {
	defer server.metrics.PeerListMsgCount.Inc()

	var addrs []peer.AddrInfo

	err := json.Unmarshal(bytes, &addrs)
	if err != nil {
		return
	}

	for _, addr := range addrs {
		if addr.ID == server.selfID {
			continue
		}

		if !server.apm.InPeers(addr.ID) {
			// NOTE: tested with remote peers
			go func(addr peer.AddrInfo) {
				ctx, cancel := context.WithTimeout(server.ctx, connectTimeout)
				defer cancel()
				server.host.Connect(ctx, addr)
			}(addr)
		}
	}
}

// process and transfer signature to the channel
func (bs *BuilderServer) optSignature(cpeer peer.ID, bytes []byte) {
	defer bs.metrics.SignatureMsgCount.Inc()

	bs.comChannels[commons.SignatureCh] <- messages.PeerMsg{
		Peer:  cpeer,
		Bytes: bytes,
	}
}

// process and transfer block key to the channel
func (bs *BuilderServer) optBlockKey(cpeer peer.ID, bytes []byte) {
	defer bs.metrics.BlockKeyMsgCount.Inc()

	bs.comChannels[commons.BlockKeyCh] <- messages.PeerMsg{
		Peer:  cpeer,
		Bytes: bytes,
	}
}

// process and transfer bundle to the channel
func (bs *BuilderServer) optBundle(cpeer peer.ID, bytes []byte) {
	defer bs.metrics.BundleMsgCount.Inc()

	bs.comChannels[commons.BundleCh] <- messages.PeerMsg{
		Peer:  cpeer,
		Bytes: bytes,
	}
}

// process and transfer preconfirmation bids to the channel
func (bs *BuilderServer) optPreconfBid(cpeer peer.ID, bytes []byte) {
	defer bs.metrics.PreconfBidMsgCount.Inc()

	bs.comChannels[commons.PreconfCh] <- messages.PeerMsg{
		Peer:  cpeer,
		Bytes: bytes,
	}
}

// process and transfer bids to the channel
func (ss *SearcherServer) optBid(cpeer peer.ID, bytes []byte) {
	defer ss.metrics.BidMsgCount.Inc()

	ss.comChannels[commons.BidCh] <- messages.PeerMsg{
		Peer:  cpeer,
		Bytes: bytes,
	}
}

// process and transfer commitments to the channel
func (ss *SearcherServer) optCommitment(cpeer peer.ID, bytes []byte) {
	defer ss.metrics.CommitmentMsgCount.Inc()

	ss.comChannels[commons.CommitmentCh] <- messages.PeerMsg{
		Peer:  cpeer,
		Bytes: bytes,
	}
}

// get self peer.ID
func (server Server) ID() peer.ID { return server.selfID }

// get approved peers on pubsub server
func (server *Server) GetApprovedPeers() []peer.ID {
	return server.apm.GetPeerIDs()
}

// get approved builders on pubsub server
func (server *Server) GetApprovedBuilderPeers() []peer.ID {
	return server.apm.GetBuilderPeerIDs()
}

// get approved searchers on pubsub server
func (server *Server) GetApprovedSearcherPeers() []peer.ID {
	return server.apm.GetSearcherPeerIDs()
}

// get approved gossip peers on pubsub server
func (server *Server) GetApprovedGossipPeers() []peer.ID {
	return server.apm.GetGossipPeerIDs()
}

// when it listens to events, it applies necessary procedures based on incoming
// or outgoing connections
func (server *Server) events(trackCh <-chan commons.ConnectionEvent, streamFuncs ...streamFunc) {

	var mu = &sync.Mutex{}
	for event := range trackCh {

		// protect events
		mu.Lock()
		eventCopy := event
		mu.Unlock()

		var stream streamFunc
		switch eventCopy.PeerType {
		case commons.Builder:
			stream = streamFuncs[0]
		case commons.Searcher:
			if server.selfType == commons.Searcher {
				//TODO unexpected situation kill connection
				// kill the Searcher connection before it reaches this point
				server.host.Network().ClosePeer(eventCopy.PeerID)
				continue

			}
			stream = streamFuncs[1]
		default:
			//NOTE unexpected situation!
			// possible race condition issue between ConnectionGater and ConnectionTracker
			server.host.Network().ClosePeer(eventCopy.PeerID)
			continue

		}

		switch eventCopy.Event {
		case commons.Connected:
			var retry = 5
			// create approve message and stream this message for new peer
			msg, err := server.omb.Approve(server.token)
			if err != nil {
				return
			}

			for i := 0; i < retry; i++ {
				err = stream(eventCopy.PeerID, msg)
				if err == nil {
					break
				}
			}

			go func() {
				checker := time.NewTicker(1000 * time.Millisecond)

				for i := 0; i < 10; i++ {
					// wait for it to join the verified peers
					<-checker.C

					if server.apm.InPeers(eventCopy.PeerID) {
						// once the peer is connected,
						// send a message to get the version information
						msg, err = server.omb.GetVersion()
						if err != nil {
							return
						}

						stream(eventCopy.PeerID, msg)
						break
					}
				}

				checker.Stop()

				if !server.apm.InPeers(eventCopy.PeerID) {
					server.host.Network().ClosePeer(eventCopy.PeerID)
					return
				}

				// generate a ping message with a unique UUID for each peer and append timestamps
				peerInfo := server.apm.GetPeerInfo(eventCopy.PeerID)
				uuidBytes, err := peerInfo.uuid.MarshalBinary()
				if err != nil {
					return
				}

				msg, err = server.omb.Ping(uuidBytes)
				if err != nil {
					return
				}

				server.apm.SetPeerInfoPingTime(eventCopy.PeerID, time.Now().UnixNano())
				stream(eventCopy.PeerID, msg)

				// retrieve the peer list from the new node and expand the network connections
				msg, err = server.omb.GetPeerList()
				if err != nil {
					return
				}

				stream(eventCopy.PeerID, msg)
				// Take into account incoming and outgoing connections to use the current
				// connected peer count within the metric.
				server.metrics.ApprovedPeerCount.Set(float64(len(server.GetApprovedPeers())))
				if server.selfType == commons.Builder {
					server.metrics.BuilderPeerCount.Set(float64(len(server.GetApprovedBuilderPeers())))
					server.metrics.SearcherPeerCount.Set(float64(len(server.GetApprovedSearcherPeers())))

				}
			}()

		case commons.Disconnected:
			server.apm.DelPeer(eventCopy.PeerID)
			// Take into account incoming and outgoing connections to use the
			// current connected peer count within the metric.
			server.metrics.ApprovedPeerCount.Set(float64(len(server.GetApprovedPeers())))
			if server.selfType == commons.Builder {
				server.metrics.BuilderPeerCount.Set(float64(len(server.GetApprovedBuilderPeers())))
				server.metrics.SearcherPeerCount.Set(float64(len(server.GetApprovedSearcherPeers())))
			}

			// delete metrics when peer is disconnected
			server.deleteMetrics(eventCopy.PeerID)
		}
	}
}

// periodically check the status of peers and update latency information
func (server *Server) latencyUpdater(streamFuncs ...streamFunc) {
	// get latency check interval
	interval := server.cfg.LatencyInterval()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		for peerID, info := range server.apm.GetPeers() {
			uuidBytes, err := info.uuid.MarshalBinary()
			if err != nil {
				continue
			}

			msg, err := server.omb.Ping(uuidBytes)
			if err != nil {
				continue
			}

			server.apm.SetPeerInfoPingTime(peerID, time.Now().UnixNano())
			var stream streamFunc
			switch info.peerType {
			case commons.Builder:
				stream = streamFuncs[0]
			case commons.Searcher:
				if server.selfType == commons.Searcher {
					//TODO unexpected situation kill connection
					// Kill the Searcher connection before it reaches this point
					server.host.Network().ClosePeer(peerID)
					continue

				}
				stream = streamFuncs[1]
			default:
				//NOTE undexpected situtation
			}

			stream(peerID, msg)
		}
	}
}

// TODO: @iowar: For a more accurate calculation, consult the team.
// the score is calculated taking into account lifetime, stake amount, and rtt values
// periodically update score information
// TODO: @iowar : New factors that will modify the scoring based on peer behavior can be included.
func (server *Server) scoreUpdater() {
	// sigmoid activation function
	sigmoid := func(x float64) float64 {
		return 1.0 / (1.0 + math.Exp(-x))
	}

	// get the score update interval from configuration
	interval := server.cfg.ScoreInterval()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		// calculate and update score for each peer
		for peerID, info := range server.apm.GetPeers() {
			// weights for each feature in the scoring formula
			w1, w2, w3 := 0.2, 0.3, 0.5

			joinDate := info.getJoinDate()
			stake := info.getStake()
			latency := info.getLatency()

			// normalize v1 by dividing milliseconds by the number of milliseconds in a 30 day
			v1 := float64(time.Now().Sub(joinDate).Milliseconds() / 3. * 864000000.)
			if v1 > 1 {
				v1 = 1.
			}

			// check peer stake if nil especially searchers
			if stake == nil {
				stake = new(big.Int)
			}

			// normalize v2 by dividing stake by 10^18
			v2 := float64(stake.Int64()) / float64(math.Pow(10, 18))
			if v2 > 1 {
				v2 = 1.
			}

			// normalize v3 by subtracting latency from 10 milliseconds and taking the min of 0
			v3 := float64(10 - latency.Milliseconds())
			if v3 < 0 {
				v3 = 0.
			}

			v3 /= 10.

			// calculate the score using the sigmoid activation function with the weighted values
			score := int(sigmoid(w1*v1+w2*v2+w3*v3) * 100.)

			// update the score for the peer
			server.apm.SetPeerInfoScore(peerID, score)

			// set score metrics
			server.metrics.ScorePeers.WithLabelValues(peerID.String()).Set(float64(score))
		}
	}
}

// delete some metrics when peer is disconnected
func (server *Server) deleteMetrics(peerID peer.ID) {
	// remove latency
	server.metrics.LatencyPeers.Delete(
		prometheus.Labels{"peer_id": peerID.String()},
	)
	// remove score
	server.metrics.ScorePeers.Delete(
		prometheus.Labels{"peer_id": peerID.String()},
	)
	// remove join date
	server.metrics.JoinDatePeers.Delete(
		prometheus.Labels{"peer_id": peerID.String()},
	)
}
