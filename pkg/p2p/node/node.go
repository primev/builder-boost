package node

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lthibault/log"

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
	CreateStream(proto string, handler func(stream network.Stream))

	// SendMsg sends a message to a peer over the specified protocol.
	SendMsg(proto protocol.ID, p peer.ID, msg string) error

	// Publish publishes a message over the topic.
	//Publish(msg []byte, err error) error
	Publish(msg message.OutboundMessage) error

	// Approve approves the node and publishes the approval message.
	Approve()

	// Close closes the node with the given reason and code.
	Close(reason string, code int)

	// Ready returns a channel that signals when the node is ready.
	Ready() <-chan struct{}

	// PreconfReader returns a channel for reading pre-confirmation bids from the node.
	PreconfReader() <-chan []byte

	// PreconfSender sends a pre-confirmation bid over the node.
	PreconfSender(preconf []byte)
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

	preconfCh chan []byte

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
	)

	// Set your own keypair
	privKey, err := crypto.UnmarshalSecp256k1PrivateKey(peerKey.D.Bytes())
	if err != nil {
		panic(err)
	}

	// gater activated intercept secured
	conngtr := newConnectionGater(rollup)

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
	connectionTracker := newConnectionTracker(logger)

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

	var (
		preconfCh = make(chan []byte)
	)

	// create pubsub server
	psio := pubsubio.New(
		ctx,
		cfg,
		logger,
		host,
		trackCh,
		token,
		am.Address,
		rollup,
		topic,
		imb,
		omb,
		preconfCh,
	)

	// fill node fields
	node := &Node{
		log:       logger,
		host:      host,
		topic:     topic,
		msgBuild:  omb,
		pubSub:    psio,
		ctx:       ctx,
		cfg:       cfg,
		token:     token,
		rollup:    rollup,
		address:   am.Address,
		stake:     stake,
		closeChan: make(chan closeSignal),
		preconfCh: preconfCh,
	}

	// start default streams
	//node.CreateStream(authProto, authStreamHandler)
	// ...

	// start peer discovery options
	go node.initDiscovery()
	// signal handler for closing
	go node.waitSignal(cancel)

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

// publish message over topic
// func (n *Node) Publish(msg []byte, err error) error {
func (n *Node) Publish(msg message.OutboundMessage) error {
	var err error

	if err != nil {
		n.log.With(log.F{
			"service":  "p2p publish",
			"err time": commons.GetNow(),
		}).Error(err)
		return err
	}

	// send message to peers
	//return n.topic.Publish(n.ctx, msg)
	return n.pubSub.Publish(msg)
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

// read preconfirmation bids from the node
func (n *Node) PreconfReader() <-chan []byte {
	return n.preconfCh
}

// publish preconfirmation bids over the node
func (n *Node) PreconfSender(preconf []byte) {
	msg, err := n.msgBuild.PreconfirmationBid(preconf)
	if err != nil {
		panic(err)
	}

	err = n.Publish(msg)
	if err != nil {
		panic(err)
	}
}
