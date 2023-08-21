package node

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lthibault/log"
	"github.com/multiformats/go-multiaddr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"github.com/primev/builder-boost/pkg/p2p/commons"
	"github.com/primev/builder-boost/pkg/p2p/commons/messages"
	"github.com/primev/builder-boost/pkg/p2p/commons/signer"
	"github.com/primev/builder-boost/pkg/p2p/config"
	"github.com/primev/builder-boost/pkg/p2p/discover"
	"github.com/primev/builder-boost/pkg/p2p/message"
	pubsubio "github.com/primev/builder-boost/pkg/p2p/pubsub"
	"github.com/primev/builder-boost/pkg/rollup"
)

var _ IBuilderNode = (*BuilderNode)(nil)
var _ ISearcherNode = (*SearcherNode)(nil)

type INode interface {
	// GetPeerID returns the peer id of the node.
	GetPeerID() peer.ID

	//GetAddrs() returns the peer addrs of the node.
	GetAddrs() []multiaddr.Multiaddr

	// GetToken returns the token of the node.
	GetToken() []byte

	// GetEthAddress returns the eth address of the node.
	GetEthAddress() common.Address

	// GetStake returns the stake amount of the node.
	GetStake() *big.Int

	// GetPeers returns the list of connected peers.
	GetPeers() peer.IDSlice

	// Close closes the node with the given reason and code.
	Close(string, int)

	// Ready returns a channel that signals when the node is ready.
	Ready() <-chan struct{}
}

type ISearcherNode interface {
	// INode extends the ISearcherNode interface and adds methods specific to node functionality.
	INode

	// GetBuilderPeersOnTopic returns the list of connected builder peers on the topic.
	GetBuilderPeersOnTopic() peer.IDSlice

	// GetApprovedBuilderPeers returns the list of approved builder peers on the psio.
	GetApprovedBuilderPeers() []peer.ID

	// Publish publishes a message over the topic.
	PublishS(message.OutboundMessage) error

	// Stream stream a message over the proto to specific peer.
	StreamS(peer.ID, message.OutboundMessage) error

	// Gossip gossip a message over the gossip proto to specific peers.
	GossipS(message.OutboundMessage) error

	// BidReader returns a channel for reading bids from the node. NOTE: currently in the testing phase
	BidReader() <-chan messages.PeerMsg

	// BidSend sends bid over the node.	NOTE: currently in the testing phase
	BidSend(commons.Broadcast, []byte) error
}

type IBuilderNode interface {
	// INode extends the IBuilderNode interface and adds methods specific to node functionality.
	INode

	// ISearcherNode extends the IBuilderNode interface and adds methods specific to searchers functionality.
	ISearcherNode

	// GetBuilderPeersOnTopic returns the list of connected builder peers on the topic.
	GetBuilderPeersOnTopic() peer.IDSlice

	// GetSearcherPeersOnTopic returns the list of connected searcher peers on the topic.
	GetSearcherPeersOnTopic() peer.IDSlice

	// GetApprovedBuilderPeers returns the list of approved builder peers on the psio.
	GetApprovedBuilderPeers() []peer.ID

	// GetApprovedBuilderPeers returns the list of approved searcher peers on the psio.
	GetApprovedSearcherPeers() []peer.ID

	// Publish publishes a message over the topic.
	PublishB(message.OutboundMessage) error

	// Stream stream a message over the proto to specific peer.
	StreamB(peer.ID, message.OutboundMessage) error

	// Gossip gossip a message over the gossip proto to specific peers.
	GossipB(message.OutboundMessage) error

	// SignatureReader returns a channel for reading signatures from the node.
	SignatureReader() <-chan messages.PeerMsg

	// BlockKeyReader returns a channel for reading block keys from the node.
	BlockKeyReader() <-chan messages.PeerMsg

	// BundleReader returns a channel for reading bundels from the node.
	BundleReader() <-chan messages.PeerMsg

	// PreconfReader returns a channel for reading pre-confirmation bids from the node.
	PreconfReader() <-chan messages.PeerMsg

	// SignatureSend sends signature over the node.
	SignatureSend(peer.ID, []byte) error

	// BlockKeySend sends block key over the node.
	BlockKeySend(peer.ID, []byte) error

	// BundleSend sends bundles over the node.
	BundleSend(commons.Broadcast, []byte) error

	// PreconfSend sends a pre-confirmation bid over the node.
	PreconfSend(commons.Broadcast, []byte) error
}

// node shutdown signal
type closeSignal struct {
	Reason string
	Code   int
}

type SearcherNode struct {
	Node
	scomChannels map[commons.ComChannels]chan messages.PeerMsg
}

type BuilderNode struct {
	Node
	SearcherNode
	bcomChannels map[commons.ComChannels]chan messages.PeerMsg
}

// specific node fields
type Node struct {
	peerType  commons.PeerType
	log       log.Logger
	host      host.Host
	topics    map[commons.Topic]*pubsub.Topic
	msgBuild  message.OutboundMsgBuilder
	psio      *pubsubio.PubSubIO
	ctx       context.Context
	cfg       *config.Config
	token     []byte
	rollup    rollup.Rollup
	address   common.Address
	stake     *big.Int
	closeChan chan closeSignal
	metrics   *metrics
	once      sync.Once
	ready     chan struct{}
}

// create builder node
func NewBuilderNode(logger log.Logger, key *ecdsa.PrivateKey, rollup rollup.Rollup, registry *prometheus.Registry) IBuilderNode {
	return newNode(logger, key, rollup, registry, commons.Builder).(IBuilderNode)
}

// create searcher node
func NewSearcherNode(logger log.Logger, key *ecdsa.PrivateKey, rollup rollup.Rollup, registry *prometheus.Registry) ISearcherNode {
	return newNode(logger, key, rollup, registry, commons.Searcher).(ISearcherNode)
}

// create p2p node
func newNode(logger log.Logger, key *ecdsa.PrivateKey, rollup rollup.Rollup, registry *prometheus.Registry, peerType commons.PeerType) interface{} {

	if logger == nil {
		switch peerType {
		case commons.Builder:
			logger = log.New().WithField("service", "p2p-builder")
		case commons.Searcher:
			logger = log.New().WithField("service", "p2p-searcher")
		default:
			return nil
		}
	}

	switch peerType {
	case commons.Builder:
	case commons.Searcher:
	default:
		logger.With(log.F{
			"caller": commons.GetCallerName(),
			"date":   commons.GetNow(),
			"mode":   peerType,
		}).Fatal("unknown peer type!")

	}

	logger.With(log.F{
		"caller": commons.GetCallerName(),
		"date":   commons.GetNow(),
		"mode":   peerType.String(),
	}).Info("p2p node initialization...")

	ctx, cancel := context.WithCancel(context.Background())

	// load config
	cfg := config.New(
		config.WithVersion("0.0.1"),
		config.WithDiscoveryInterval(30*time.Minute),
		config.WithLatencyInterval(time.Hour*1),
		config.WithScoreInterval(time.Minute*1),
		config.WithMinimalStake(big.NewInt(1)),
		config.WithMetricsNamespace("primev"),
		config.WithMetricsPort(8081),
		config.WithMetricsRoute("/metrics_p2p"),
	)

	// Set your own keypair
	privKey, err := crypto.UnmarshalSecp256k1PrivateKey(key.D.Bytes())
	if err != nil {
		logger.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p keypair",
		}).Fatal(err)
	}

	// create metrics
	// @iowar: If system metrics need to be combined, use the main system's registry variable
	var newRegistry bool
	if registry == nil {
		registry = prometheus.NewRegistry()
		newRegistry = true
	} else {
		newRegistry = false
	}
	metrics := newMetrics(registry, cfg.MetricsNamespace())

	// create blocker
	blocker := newBlocker(metrics, logger)

	// gater activated intercept secured
	conngtr := newConnectionGater(metrics, rollup, blocker, logger, cfg.MinimalStake(), peerType)

	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		logger.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p connmngr",
		}).Fatal(err)
	}
	host, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(privKey),
		// Connection gater
		libp2p.ConnectionGater(conngtr),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0", // regular tcp connections
			//"/ip4/0.0.0.0/tcp/0/ws", // websocket endpoint
			//"/ip4/0.0.0.0/udp/0/quic", // a UDP endpoint for the QUIC transport
		),
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connections
		libp2p.Security(noise.ID, noise.New),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		// Let this host use the DHT to find other hosts
		//libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
		//	idht, err = dht.New(ctx, h)
		//	return idht, err
		//}),
		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		libp2p.EnableNATService(),
	)
	if err != nil {
		logger.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p init",
		}).Fatal(err)
	}

	for _, addr := range host.Addrs() {
		logger.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p host",
			"addr":    addr,
			"host":    host.ID().Pretty(),
		}).Info("host address")
	}

	// create a connectionTracker for connection tracking
	connectionTracker := newConnectionTracker(conngtr, metrics, logger)

	// listen to network events
	host.Network().Notify(&network.NotifyBundle{
		ConnectedF:    connectionTracker.handleConnected,
		DisconnectedF: connectionTracker.handleDisconnected,
	})

	connTrackCh := connectionTracker.trackCh

	// create a new PubSub instance
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		logger.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p create pubsub server",
		}).Fatal(err)
	}

	// NOTE: direct publish operations are deprecated
	createTopics := func() map[commons.Topic]*pubsub.Topic {
		switch peerType {
		case commons.Builder:
			topic_b2b, err := ps.Join(cfg.PubSubTopicB2B())
			if err != nil {
				logger.With(log.F{
					"caller":  commons.GetCallerName(),
					"date":    commons.GetNow(),
					"service": "p2p pubsub b2b topic join",
				}).Fatal(err)
			}

			topic_b2s, err := ps.Join(cfg.PubSubTopicB2S())
			if err != nil {
				logger.With(log.F{
					"caller":  commons.GetCallerName(),
					"date":    commons.GetNow(),
					"service": "p2p pubsub b2s topic join",
				}).Fatal(err)
			}

			return map[commons.Topic]*pubsub.Topic{
				commons.TopicB2B: topic_b2b,
				commons.TopicB2S: topic_b2s,
			}

		case commons.Searcher:
			topic_b2s, err := ps.Join(cfg.PubSubTopicB2S())
			if err != nil {
				logger.With(log.F{
					"caller":  commons.GetCallerName(),
					"date":    commons.GetNow(),
					"service": "p2p pubsub b2s topic join",
				}).Fatal(err)
			}

			return map[commons.Topic]*pubsub.Topic{
				commons.TopicB2S: topic_b2s,
			}
		default:
		}

		return nil
	}

	// create inbound and outbound message builders
	imb, omb := message.NewInboundBuilder(), message.NewOutboundBuilder()

	token, err := generateToken(peerType, host.ID(), key)
	if err != nil {
		logger.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p generate builder token",
		}).Fatal(err)
	}

	// eth address
	address := commons.GetAddressFromPrivateKey(key)

	// get stake amount of builder
	var (
		stake = big.NewInt(0)
	)
	switch peerType {
	case commons.Builder:
		stake, err := rollup.GetMinimalStake(address)
		if err != nil {
			logger.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p rollup",
			}).Fatal(err)
		}

		// if there is not enough stake, close node
		if stake.Cmp(cfg.MinimalStake()) < 0 {
			err = errors.New("not enough stake")
			logger.With(log.F{
				"service":  "p2p minimal stake",
				"log time": commons.GetNow(),
			}).Fatal(err)
		}
	case commons.Searcher:
		//TODO no scenario has been designed for Searchers yet, follow the team opinions
	}

	var topics = createTopics()

	var (
		scomChannels map[commons.ComChannels]chan messages.PeerMsg
		bcomChannels map[commons.ComChannels]chan messages.PeerMsg
	)

	switch peerType {
	case commons.Builder:
		scomChannels = map[commons.ComChannels]chan messages.PeerMsg{
			commons.BidCh: make(chan messages.PeerMsg),
		}

		bcomChannels = map[commons.ComChannels]chan messages.PeerMsg{
			commons.SignatureCh: make(chan messages.PeerMsg),
			commons.BlockKeyCh:  make(chan messages.PeerMsg),
			commons.BundleCh:    make(chan messages.PeerMsg),
			commons.PreconfCh:   make(chan messages.PeerMsg),
		}

	case commons.Searcher:
		scomChannels = map[commons.ComChannels]chan messages.PeerMsg{
			commons.BidCh: make(chan messages.PeerMsg),
		}
	}

	comChannels := mergeChannelMaps(scomChannels, bcomChannels)

	// create pubsub server
	psio, err := pubsubio.New(
		ctx,
		cfg,
		logger,
		peerType,
		host,
		registry,
		connTrackCh,
		token,
		address,
		rollup,
		topics,
		imb,
		omb,
		comChannels,
	)
	if err != nil {
		logger.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "pubsubio server initialization",
		}).Fatal(err)
	}

	// fill node fields
	node := &Node{
		peerType:  peerType,
		log:       logger,
		host:      host,
		topics:    topics,
		msgBuild:  omb,
		psio:      psio,
		ctx:       ctx,
		cfg:       cfg,
		token:     token,
		rollup:    rollup,
		address:   address,
		stake:     stake,
		closeChan: make(chan closeSignal),
		metrics:   metrics,
	}

	searcherNode := &SearcherNode{
		Node:         *node,
		scomChannels: scomChannels,
	}

	builderNode := &BuilderNode{
		Node:         *node,
		SearcherNode: *searcherNode,
		bcomChannels: bcomChannels,
	}

	// start default streams
	//node.CreateStream(authProto, authStreamHandler)
	// ...

	// start peer discovery options
	go node.initDiscovery()
	// signal handler for closing
	go node.waitSignal(cancel)

	// start p2p metrics handler
	if newRegistry {
		go node.initMetrics(registry)
	}

	// TODO the temporary hold will be lifted once the discovery options are provided
	//time.Sleep(time.Second * 3)

	switch peerType {
	case commons.Builder:
		builderNode.setReady()
		return builderNode
	case commons.Searcher:
		searcherNode.setReady()
		return searcherNode
	}
	return nil
}

// get peer id
func (n *Node) GetPeerID() peer.ID {
	return n.host.ID()
}

// get peer addrs
func (n *Node) GetAddrs() []multiaddr.Multiaddr {
	return n.host.Addrs()
}

// get node signed auth token
func (n *Node) GetToken() []byte {
	return n.token
}

// get node eth address
func (n *Node) GetEthAddress() common.Address {
	return n.address
}

// get stake amount (rollup builder: on|rollup searcher=(off))
// TODO: rollup improvement
func (n *Node) GetStake() *big.Int {
	return n.stake
}

// get connected peer list
func (n *Node) GetPeers() peer.IDSlice {
	return n.host.Peerstore().Peers()
}

// get connected builder peer list on topic
func (bn *BuilderNode) GetBuilderPeersOnTopic() peer.IDSlice {
	return bn.topics[commons.TopicB2B].ListPeers()
}

// get connected searcher peer list on topic
func (bn *BuilderNode) GetSearcherPeersOnTopic() peer.IDSlice {
	return bn.topics[commons.TopicB2S].ListPeers()
}

// get connected builder peer list on topic
func (sn *SearcherNode) GetBuilderPeersOnTopic() peer.IDSlice {
	return sn.topics[commons.TopicB2S].ListPeers()
}

// get approved builder peers
func (bn *BuilderNode) GetApprovedBuilderPeers() []peer.ID {
	return bn.psio.BuilderServer.GetApprovedBuilderPeers()
}

// get approved searcher peers
func (bn *BuilderNode) GetApprovedSearcherPeers() []peer.ID {
	return bn.psio.BuilderServer.GetApprovedBuilderPeers()
}

// get approved builder peers
func (sn *SearcherNode) GetApprovedBuilderPeers() []peer.ID {
	return sn.psio.SearcherServer.GetApprovedBuilderPeers()
}

// create new stream proto
func (n *Node) CreateStream(proto string, handler func(stream network.Stream)) {
	n.host.SetStreamHandler(protocol.ID(proto), handler)
}

// send message to peer over given protocol
// TODO this method is used for testing purposes and will be removed in the future
func (n *Node) SendMsg(proto protocol.ID, p peer.ID, msg string) error {
	s, err := n.host.NewStream(n.ctx, p, proto)
	if err != nil {
		return err
	}

	defer s.Close()

	w := bufio.NewWriter(s)
	l, err := w.WriteString(msg)
	if l != len(msg) {
		return fmt.Errorf("expected to write %d bytes, wrote %d", len(msg), l)
	}
	if err != nil {
		return err
	}
	if err = w.Flush(); err != nil {
		return err
	}

	return nil
}

// publish message over pubsub topic
func (bn *BuilderNode) PublishB(msg message.OutboundMessage) error {
	// send message to all approved peers
	return bn.psio.BuilderServer.Publish(msg)
}

// stream message over pubsub stream proto
func (bn *BuilderNode) StreamB(peerID peer.ID, msg message.OutboundMessage) error {
	// stream message to specific peer
	return bn.psio.BuilderServer.Stream(peerID, msg)
}

// gossip message over pubsub gossip proto
func (bn *BuilderNode) GossipB(msg message.OutboundMessage) error {
	// gossip message to specific peers
	return bn.psio.BuilderServer.Gossip(msg)
}

// publish message over pubsub topic
func (sn *SearcherNode) PublishS(msg message.OutboundMessage) error {
	// send message to all approved peers
	return sn.psio.SearcherServer.Publish(msg)
}

// stream message over pubsub stream proto
func (sn *SearcherNode) StreamS(peerID peer.ID, msg message.OutboundMessage) error {
	// stream message to specific peer
	return sn.psio.SearcherServer.Stream(peerID, msg)
}

// gossip message over pubsub gossip proto
func (sn *SearcherNode) GossipS(msg message.OutboundMessage) error {
	// gossip message to specific peers
	return sn.psio.SearcherServer.Gossip(msg)
}

// send close signal to node
func (n *Node) Close(reason string, code int) {
	signal := closeSignal{
		Reason: reason,
		Code:   code,
	}

	n.closeChan <- signal
	close(n.closeChan)
}

// initial discovery options
// mdns, bootstrapt, dht etc.
func (n *Node) initDiscovery() {
	discovery := discover.NewDiscovery(n.cfg, n.host, n.ctx, n.log)

	// setup local mDNS discovery
	if err := discovery.StartMdnsDiscovery(); err != nil {
		n.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p mdns discovery",
		}).Warn(err)
	}

	// connect default nodes
	discovery.ConnectBootstrap()

	// It has a weak infrastructure in libp2p
	discovery.StartDhtRouting()
}

// trigger with signal and do tasks safely
func (n *Node) waitSignal(cancel context.CancelFunc) {
	signal := <-n.closeChan

	// ps.Leave(config.Topic)

	// close topics
	for _, topic := range n.topics {
		topic.Close()
	}

	// close host
	n.host.Close()
	// context cancel
	cancel()

	n.log.With(log.F{
		"caller":  commons.GetCallerName(),
		"date":    commons.GetNow(),
		"service": "p2p node",
		"reason":  signal.Reason,
		"code":    signal.Code,
	}).Warn("node is being turned off...")
}

func (n *Node) Ready() <-chan struct{} {
	n.once.Do(func() {
		n.ready = make(chan struct{})
	})
	return n.ready
}

func (n *Node) initMetrics(registry *prometheus.Registry) {
	var (
		port  = n.cfg.MetricsPort()
		route = n.cfg.MetricsRoute()
	)
	n.log.With(log.F{
		"caller":  commons.GetCallerName(),
		"date":    commons.GetNow(),
		"service": "p2p metrics",
		"port":    port,
		"route":   route,
	}).Info("p2p metrics started...")

	promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.Handle(route, promHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (n *Node) setReady() {
	select {
	case <-n.Ready():
	default:
		close(n.ready)
	}
}

// read signatures from the node
func (bn *BuilderNode) SignatureReader() <-chan messages.PeerMsg {
	return bn.bcomChannels[commons.SignatureCh]
}

// read block keys from the node
func (bn *BuilderNode) BlockKeyReader() <-chan messages.PeerMsg {
	return bn.bcomChannels[commons.BlockKeyCh]
}

// read bundles from the node
func (bn *BuilderNode) BundleReader() <-chan messages.PeerMsg {
	return bn.bcomChannels[commons.BundleCh]
}

// read preconfirmation bids from the node
func (bn *BuilderNode) PreconfReader() <-chan messages.PeerMsg {
	return bn.bcomChannels[commons.PreconfCh]
}

// read bids from the node
func (sn *SearcherNode) BidReader() <-chan messages.PeerMsg {
	return sn.scomChannels[commons.BidCh]
}

// stream signature over the node for specific peer
func (bn *BuilderNode) SignatureSend(peer peer.ID, sig []byte) error {
	msg, err := bn.msgBuild.Signature(sig)
	if err != nil {
		return err
	}

	err = bn.StreamB(peer, msg)
	if err != nil {
		return err
	}

	return nil
}

// stream block key over the node for specific peer
func (bn *BuilderNode) BlockKeySend(peer peer.ID, sig []byte) error {
	msg, err := bn.msgBuild.BlockKey(sig)
	if err != nil {
		return err
	}

	err = bn.StreamB(peer, msg)
	if err != nil {
		return err
	}

	return nil
}

// gossip or publish bundles over the node
func (bn *BuilderNode) BundleSend(broadcastType commons.Broadcast, bundle []byte) error {
	msg, err := bn.msgBuild.Bundle(bundle)
	if err != nil {
		return err
	}

	if broadcastType == commons.Publish {
		err = bn.PublishB(msg)
		if err != nil {
			return err
		}
	} else if broadcastType == commons.Gossip {
		err = bn.GossipB(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// gossip or publish preconfirmation bids over the node
func (bn *BuilderNode) PreconfSend(broadcastType commons.Broadcast, preconf []byte) error {
	msg, err := bn.msgBuild.PreconfBid(preconf)
	if err != nil {
		return err
	}

	if broadcastType == commons.Publish {
		err = bn.PublishB(msg)
		if err != nil {
			return err
		}
	} else if broadcastType == commons.Gossip {
		err = bn.GossipB(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// gossip or publish bids over the node
func (sn *SearcherNode) BidSend(broadcastType commons.Broadcast, bid []byte) error {
	msg, err := sn.msgBuild.Bid(bid)
	if err != nil {
		return err
	}

	if broadcastType == commons.Publish {
		err = sn.PublishS(msg)
		if err != nil {
			return err
		}
	} else if broadcastType == commons.Gossip {
		err = sn.GossipS(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// generate token for peer
func generateToken(peerType commons.PeerType, peerID peer.ID, key *ecdsa.PrivateKey) ([]byte, error) {
	newSigner := signer.New()

	// eth address
	address := commons.GetAddressFromPrivateKey(key)

	am := messages.ApproveMsg{
		PeerType: peerType,
		Peer:     peerID,
		Address:  address,
	}

	msgBytes := am.GetUnsignedMessage()
	sig, err := newSigner.Sign(key, msgBytes)
	if err != nil {
		return nil, err
	}

	am.Sig = sig
	token, err := json.Marshal(am)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// bringing together communication channels through merging
func mergeChannelMaps(maps ...map[commons.ComChannels]chan messages.PeerMsg) map[commons.ComChannels]chan messages.PeerMsg {
	combinedMap := make(map[commons.ComChannels]chan messages.PeerMsg)

	for _, m := range maps {
		for k, v := range m {
			combinedMap[k] = v
		}
	}

	return combinedMap
}
