package node

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lthibault/log"
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

// BoostNode interface defines the functionality of a P2P node.
type BoostNode interface {
	// GetToken returns the token of the node.
	GetToken() []byte

	// GetAddress returns the address of the node.
	GetAddress() common.Address

	// GetStake returns the stake amount of the node.
	GetStake() *big.Int

	// GetPeers returns the list of connected peers.
	GetPeers() peer.IDSlice

	// GetPeersOnTopic returns the list of connected peers on the topic.
	GetPeersOnTopic() peer.IDSlice

	// GetApprovedPeers returns the list of approved peers.
	GetApprovedPeers() []peer.ID

	// CreateStream creates a new stream with the given protocol and handler function.
	CreateStream(string, func(stream network.Stream))

	// SendMsg sends a message to a peer over the specified protocol.
	SendMsg(protocol.ID, peer.ID, string) error

	// Publish publishes a message over the topic.
	Publish(message.OutboundMessage) error

	// Stream stream a message over the proto to specific peer.
	Stream(peer.ID, message.OutboundMessage) error

	// Gossip gossip a message over the gossip proto to specific peers.
	Gossip(message.OutboundMessage) error

	// Approve approves the node and publishes the approval message.
	Approve()

	// Close closes the node with the given reason and code.
	Close(string, int)

	// Ready returns a channel that signals when the node is ready.
	Ready() <-chan struct{}

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

// specific node fields
type Node struct {
	log      log.Logger
	host     host.Host
	topic    *pubsub.Topic
	msgBuild message.OutboundMsgBuilder
	pubSub   *pubsubio.Server

	ctx       context.Context
	cfg       *config.Config
	token     []byte
	rollup    rollup.Rollup
	address   common.Address
	stake     *big.Int
	closeChan chan closeSignal

	signatureCh chan messages.PeerMsg
	blockKeyCh  chan messages.PeerMsg
	bundleCh    chan messages.PeerMsg
	preconfCh   chan messages.PeerMsg

	metrics *metrics

	once  sync.Once
	ready chan struct{}
}

// create p2p node
func CreateNode(logger log.Logger, peerKey *ecdsa.PrivateKey, rollup rollup.Rollup) BoostNode {
	if logger == nil {
		logger = log.New().WithField("service", "p2p")
	}

	logger.With(log.F{
		"service":    "p2p createnode",
		"start time": commons.GetNow(),
	}).Info("starting node...")

	ctx, cancel := context.WithCancel(context.Background())

	// load config
	cfg := config.New(
		config.WithVersion("0.0.2"),
		config.WithDiscoveryInterval(30*time.Minute),
		config.WithLatencyInterval(time.Hour*1),
		config.WithScoreInterval(time.Minute*10),
		config.WithMinimalStake(big.NewInt(1)),
		config.WithMetricsNamespace("primev"),
		config.WithMetricsPort(8081),
	)

	// Set your own keypair
	privKey, err := crypto.UnmarshalSecp256k1PrivateKey(peerKey.D.Bytes())
	if err != nil {
		panic(err)
	}

	// create metrics
	registry := prometheus.NewRegistry()
	metrics := newMetrics(registry, cfg.MetricsNamespace())

	// gater activated intercept secured
	conngtr := newConnectionGater(metrics, rollup, cfg.MinimalStake())

	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	for _, addr := range host.Addrs() {
		logger.With(log.F{
			"service":  "p2p host",
			"addr":     addr,
			"host":     host.ID().Pretty(),
			"log time": commons.GetNow(),
		}).Info("host address")
	}

	// create a connectionTracker for connection tracking
	connectionTracker := newConnectionTracker(metrics, logger)

	// listen to network events
	host.Network().Notify(&network.NotifyBundle{
		ConnectedF:    connectionTracker.handleConnected,
		DisconnectedF: connectionTracker.handleDisconnected,
	})

	trackCh := connectionTracker.trackCh

	// create a new PubSub instance
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}

	// direct publish operations are deprecated
	// join the topic as a publisher
	topic, err := ps.Join(cfg.PubSubTopic())
	if err != nil {
		panic(err)
	}

	// create inbound and outbound message builders
	imb, omb := message.NewInboundBuilder(), message.NewOutboundBuilder()

	// generate token
	newSigner := signer.New()

	am := messages.ApproveMsg{
		Peer:    host.ID(),
		Address: commons.GetAddressFromPrivateKey(peerKey),
	}

	msgBytes := am.GetUnsignedMessage()
	sig, err := newSigner.Sign(peerKey, msgBytes)
	if err != nil {
		logger.With(log.F{
			"service":  "signer",
			"err time": commons.GetNow(),
		}).Error(err)
		panic(err)
	}

	am.Sig = sig
	token, err := json.Marshal(am)
	if err != nil {
		panic(err)
	}

	// get stake amount
	stake, err := rollup.GetMinimalStake(am.Address)
	if err != nil {
		panic(err)
	}

	// create ofx channels
	var (
		signatureCh = make(chan messages.PeerMsg)
		blockKeyCh  = make(chan messages.PeerMsg)
		bundleCh    = make(chan messages.PeerMsg)
		preconfCh   = make(chan messages.PeerMsg)
	)

	// create pubsub server
	psio := pubsubio.New(
		ctx,
		cfg,
		logger,
		host,
		registry,
		trackCh,
		token,
		am.Address,
		rollup,
		topic,
		imb,
		omb,
		signatureCh,
		blockKeyCh,
		bundleCh,
		preconfCh,
	)

	// fill node fields
	node := &Node{
		log:         logger,
		host:        host,
		topic:       topic,
		msgBuild:    omb,
		pubSub:      psio,
		ctx:         ctx,
		cfg:         cfg,
		token:       token,
		rollup:      rollup,
		address:     am.Address,
		stake:       stake,
		closeChan:   make(chan closeSignal),
		signatureCh: signatureCh,
		blockKeyCh:  blockKeyCh,
		bundleCh:    bundleCh,
		preconfCh:   preconfCh,
		metrics:     metrics,
	}

	// start default streams
	//node.CreateStream(authProto, authStreamHandler)
	// ...

	// start peer discovery options
	go node.initDiscovery()
	// signal handler for closing
	go node.waitSignal(cancel)

	// start p2p metrics handler
	go func() {
		promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		http.Handle("/metrics_p2p", promHandler)
		http.ListenAndServe(fmt.Sprintf(":%d", cfg.MetricsPort()), nil)
	}()

	// TODO the temporary hold will be lifted once the discovery options are provided
	time.Sleep(time.Second * 3)
	node.setReady()

	return node
}

func (n *Node) GetToken() []byte {
	return n.token
}

func (n *Node) GetAddress() common.Address {
	return n.address
}

func (n *Node) GetStake() *big.Int {
	return n.stake
}

// get connected peer list
func (n *Node) GetPeers() peer.IDSlice {
	return n.host.Peerstore().Peers()
}

// get connected peer list on topic
func (n *Node) GetPeersOnTopic() peer.IDSlice {
	return n.topic.ListPeers()
}

// get approved peers
func (n *Node) GetApprovedPeers() []peer.ID {
	return n.pubSub.GetApprovedPeers()
}

// create new stream proto
func (n *Node) CreateStream(proto string, handler func(stream network.Stream)) {
	n.host.SetStreamHandler(protocol.ID(proto), handler)
}

// send message to peer over given protocol
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
func (n *Node) Publish(msg message.OutboundMessage) error {
	// send message to all approved peers
	return n.pubSub.Publish(msg)
}

// stream message over pubsub stream proto
func (n *Node) Stream(peerID peer.ID, msg message.OutboundMessage) error {
	// stream message to specific peer
	return n.pubSub.Stream(peerID, msg)
}

// gossip message over pubsub gossip proto
func (n *Node) Gossip(msg message.OutboundMessage) error {
	// gossip message to specific peers
	return n.pubSub.Gossip(msg)
}

// approve over node
func (n *Node) Approve() {
	msg, err := n.msgBuild.Approve(n.GetToken())
	if err != nil {
		panic(err)
	}

	//err = n.Publish(msg.MarshalJSON())
	err = n.Publish(msg)
	if err != nil {
		panic(err)
	}
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
		panic(err)
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

	// close topic
	n.topic.Close()
	// close host
	n.host.Close()
	// context cancel
	cancel()

	n.log.With(log.F{
		"service":    "p2p node",
		"reason":     signal.Reason,
		"code":       signal.Code,
		"start time": commons.GetNow(),
	}).Info("node is being turned off...")
}

func (n *Node) Ready() <-chan struct{} {
	n.once.Do(func() {
		n.ready = make(chan struct{})
	})
	return n.ready
}

func (n *Node) setReady() {
	select {
	case <-n.Ready():
	default:
		close(n.ready)
	}
}

// read signatures from the node
func (n *Node) SignatureReader() <-chan messages.PeerMsg {
	return n.signatureCh
}

// read block keys from the node
func (n *Node) BlockKeyReader() <-chan messages.PeerMsg {
	return n.blockKeyCh
}

// read bundles from the node
func (n *Node) BundleReader() <-chan messages.PeerMsg {
	return n.bundleCh
}

// read preconfirmation bids from the node
func (n *Node) PreconfReader() <-chan messages.PeerMsg {
	return n.preconfCh
}

// stream signature over the node for specific peer
func (n *Node) SignatureSend(peer peer.ID, sig []byte) error {
	msg, err := n.msgBuild.Signature(sig)
	if err != nil {
		return err
	}

	err = n.Stream(peer, msg)
	if err != nil {
		return err
	}

	return nil
}

// stream block key over the node for specific peer
func (n *Node) BlockKeySend(peer peer.ID, sig []byte) error {
	msg, err := n.msgBuild.BlockKey(sig)
	if err != nil {
		return err
	}

	err = n.Stream(peer, msg)
	if err != nil {
		return err
	}

	return nil
}

// gossip or publish bundles over the node
func (n *Node) BundleSend(broadcastType commons.Broadcast, bundle []byte) error {
	msg, err := n.msgBuild.Bundle(bundle)
	if err != nil {
		return err
	}

	if broadcastType == commons.Publish {
		err = n.Publish(msg)
		if err != nil {
			return err
		}
	} else if broadcastType == commons.Gossip {
		err = n.Gossip(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// gossip or publish preconfirmation bids over the node
func (n *Node) PreconfSend(broadcastType commons.Broadcast, preconf []byte) error {
	msg, err := n.msgBuild.PreconfBid(preconf)
	if err != nil {
		return err
	}

	if broadcastType == commons.Publish {
		err = n.Publish(msg)
		if err != nil {
			return err
		}
	} else if broadcastType == commons.Gossip {
		err = n.Gossip(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
