package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"math"
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
)

type Server struct {
	ctx     context.Context
	cfg     *config.Config
	log     log.Logger
	self    peer.ID
	host    host.Host
	trackCh chan commons.ConnectionEvent
	token   []byte
	address common.Address
	rollup  rollup.Rollup
	topic   *pubsub.Topic
	imb     message.InboundMsgBuilder
	omb     message.OutboundMsgBuilder
	apm     *approvedPeersMap
	psp     stream.Stream

	preconfCh chan []byte

	ready chan struct{}
}

func New(
	ctx context.Context,
	cfg *config.Config,
	log log.Logger,
	host host.Host,
	trackCh chan commons.ConnectionEvent,
	token []byte,
	address common.Address,
	rollup rollup.Rollup,
	topic *pubsub.Topic,
	imb message.InboundMsgBuilder,
	omb message.OutboundMsgBuilder,
	preconfCh chan []byte,
) *Server {
	pss := new(Server)
	apm := newApprovedPeersMap()

	pss = &Server{
		ctx:       ctx,
		cfg:       cfg,
		log:       log,
		self:      host.ID(),
		host:      host,
		trackCh:   trackCh,
		token:     token,
		address:   address,
		rollup:    rollup,
		topic:     topic,
		imb:       imb,
		omb:       omb,
		apm:       apm,
		preconfCh: preconfCh,
	}

	pss.ready = make(chan struct{})

	// event tracking
	go pss.events(trackCh)
	go pss.latencyUpdater()

	// create peer stream protocol
	pss.psp = stream.New(
		host,
		protocol.ID(cfg.PeerStreamProto()),
		pss.peerStreamHandler,
	)

	var once sync.Once
	go pss.baseProtocol(once)

	<-pss.ready
	close(pss.ready)

	return pss
}

// base pubsub procotol
func (pss *Server) baseProtocol(once sync.Once) {
	pss.log.With(log.F{
		"service":    "p2p pubsub",
		"start time": commons.GetNow(),
	}).Info("starting pubsub protocol...")

	sub, err := pss.topic.Subscribe()
	if err != nil {
		panic(err)
	}

	defer sub.Cancel()
	for {
		select {
		case <-pss.ctx.Done():
			pss.log.With(log.F{
				"service":   "p2p pubsub",
				"stop time": commons.GetNow(),
			}).Info("stopped pubsub protocol...")

			return
		default:
		}

		once.Do(func() {
			pss.ready <- struct{}{}
		})

		msg, err := sub.Next(pss.ctx)
		if err != nil {
			pss.log.With(log.F{
				"service":  "p2p pubsub protocol reader",
				"err time": commons.GetNow(),
			}).Error(err)
			continue
		}

		// check self peer.ID
		if strings.Compare(pss.self.String(), msg.ReceivedFrom.String()) == 0 {
			continue
		}

		inMsg, err := pss.imb.Parse(msg.ReceivedFrom, msg.Data)
		if err != nil {
			pss.log.With(log.F{
				"service":  "p2p pubsub protocol",
				"err time": commons.GetNow(),
			}).Error(err)
			continue
		}

		// if not in approved peers try to authenticate
		if !pss.apm.InPeers(msg.ReceivedFrom) {
			go func() {
				pss.log.With(log.F{
					"service":  "p2p pubsub",
					"op":       inMsg.Op(),
					"peer":     inMsg.Peer(),
					"msg time": commons.GetNow(),
				}).Info("unverified peer message")

				switch inMsg.Op() {
				case message.Approve:
					pss.optApprove(inMsg.Peer(), inMsg.Bytes(), true)
				default:
					pss.log.With(log.F{
						"service":  "p2p pubsub",
						"op":       inMsg.Op(),
						"peer":     inMsg.Peer(),
						"msg time": commons.GetNow(),
					}).Info("unknown approve option!")
				}
			}()
			continue
		}

		go func() {
			pss.log.With(log.F{
				"service":  "p2p pubsub",
				"op":       inMsg.Op(),
				"peer":     inMsg.Peer(),
				"msg time": commons.GetNow(),
			}).Info("verified peer message")

			switch inMsg.Op() {
			// pass auth option in this side for now
			case message.Approve:
			// create pong message and publish to show you're alive
			case message.Ping:
				pss.optPing(inMsg.Peer(), inMsg.Bytes())
			// it can be use validate peer is alive
			case message.Pong:
				pss.optPong(inMsg.Peer(), inMsg.Bytes())
			// publish version
			case message.GetVersion:
				pss.optGetVersion(inMsg.Peer())
			// store peers version
			case message.Version:
				pss.optVersion(inMsg.Peer(), inMsg.Bytes())
			// publish approved peers
			case message.GetPeerList:
				pss.optGetPeerList(inMsg.Peer())
			// it can be use for connect to peers
			case message.PeerList:
				pss.optPeerList(inMsg.Peer(), inMsg.Bytes())

			case message.PreconfBid:
				pss.optPreconfBid(inMsg.Peer(), inMsg.Bytes())
			default:
				pss.log.With(log.F{
					"service":  "p2p pubsub",
					"op":       inMsg.Op(),
					"peer":     inMsg.Peer(),
					"msg time": commons.GetNow(),
				}).Info("unknown option!")
			}
		}()
	}
}

// publish message on topic
func (pss *Server) Publish(msg message.OutboundMessage) error {
	msgBytes, err := msg.MarshalJSON()
	if err != nil {
		return err
	}

	return pss.topic.Publish(pss.ctx, msgBytes)
}

// stream message for specific peer
func (pss *Server) Stream(peerID peer.ID, msg message.OutboundMessage) error {
	msgBytes, err := msg.MarshalJSON()
	if err != nil {
		return err
	}

	return pss.psp.Send(peerID, msgBytes)
}

// stream message for gossip peers
func (pss *Server) Gossip(msg message.OutboundMessage) error {
	msgBytes, err := msg.MarshalJSON()
	if err != nil {
		return err
	}

	// TODO use optimization to fast data distribution
	peers := pss.apm.GetGossipPeers()

	for peerID, _ := range peers {
		err = pss.psp.Send(peerID, msgBytes)
		if err != nil {
			pss.log.With(log.F{
				"service":    "p2p pubsub gossip",
				"close time": commons.GetNow(),
				"peer id":    peerID,
			}).Error(err)
		}
	}

	return nil
}

func (pss *Server) optApprove(cpeer peer.ID, bytes []byte, sendauth bool) {
	var am = new(messages.ApproveMsg)
	err := json.Unmarshal(bytes, &am)

	newSigner := signer.New()
	valid, address, err := newSigner.Verify(am.Sig, am.GetUnsignedMessage())
	if err != nil {
		pss.log.With(log.F{
			"service":    "p2p pubsub approve",
			"close time": commons.GetNow(),
		}).Error(err)

		// terminate the unexpected connection
		pss.host.Network().ClosePeer(cpeer)
		return
	}

	// verify signature
	if !valid {
		pss.log.With(log.F{
			"service":    "p2p pubsub approve",
			"close time": commons.GetNow(),
		}).Error(errors.New("not valid signature"))

		// terminate the unexpected connection
		pss.host.Network().ClosePeer(cpeer)
		return

	}

	// verify the correct builder eth address packaging
	if !commons.BytesCompare(address.Bytes(), am.Address.Bytes()) {
		pss.log.With(log.F{
			"service":    "p2p pubsub approve",
			"close time": commons.GetNow(),
		}).Error(errors.New("wrong adress packaging"))

		// terminate the unexpected connection
		pss.host.Network().ClosePeer(cpeer)
		return
	}

	// terminate clone builder connections
	if commons.BytesCompare(pss.address.Bytes(), am.Address.Bytes()) {
		pss.log.With(log.F{
			"service":    "p2p pubsub approve",
			"close time": commons.GetNow(),
		}).Error(errors.New("clone builder detected"))

		// terminate the unexpected connection
		pss.host.Network().ClosePeer(cpeer)
		return
	}

	// verify the correct builder peer address packaging
	if strings.Compare(cpeer.String(), am.Peer.String()) != 0 {
		pss.log.With(log.F{
			"service":    "p2p pubsub approve",
			"close time": commons.GetNow(),
		}).Error(errors.New("wrong peer packaging"))

		// terminate the unexpected connection
		pss.host.Network().ClosePeer(cpeer)
		return
	}

	// get stake amount
	stake, err := pss.rollup.GetMinimalStake(address)
	if err != nil {
		return
	}

	// check builder stake amount
	if stake.Cmp(pss.cfg.MinimalStake()) >= 0 {
		pss.apm.AddPeer(cpeer)
		for p, i := range pss.apm.GetPeers() {
			switch pss.host.Network().Connectedness(p) {
			case network.NotConnected:
				pss.apm.DelPeer(p)
			case network.CannotConnect:
				pss.apm.DelPeer(p)
			}

			// close the old connection if builder tries to establish a new one without closing it
			if commons.BytesCompare(i.address.Bytes(), am.Address.Bytes()) {
				pss.apm.DelPeer(p)
				if pss.host.Network().Connectedness(p) == network.Connected {
					pss.host.Network().ClosePeer(p)
				}
			}
		}

		addrInfo := pss.host.Peerstore().PeerInfo(cpeer)

		pss.apm.SetPeerInfoStart(cpeer, time.Now())
		pss.apm.SetPeerInfoAddress(cpeer, address)
		pss.apm.SetPeerInfoStake(cpeer, stake)
		pss.apm.SetPeerInfoAddrs(cpeer, addrInfo.Addrs)
		pss.apm.SetPeerInfoUUID(cpeer, uuid.New())
		// temporary gossip
		pss.apm.SetPeerInfoGossip(cpeer, true)

		if sendauth {
			// create approve message and stream this message for new peer
			msg, err := pss.omb.Approve(pss.token)
			if err != nil {
				return
			}

			pss.Stream(cpeer, msg)
		}
	} else {
		// terminate the unexpected connection
		pss.host.Network().ClosePeer(cpeer)

		pss.log.With(log.F{
			"service":    "p2p pubsub approve",
			"peer":       cpeer,
			"close time": commons.GetNow(),
		}).Info("not enough stake")
	}
}

// Create a 'pong' message and stream it to the peer that received the 'ping' message.
func (pss *Server) optPing(cpeer peer.ID, uuidBytes []byte) {
	msg, err := pss.omb.Pong(uuidBytes)
	if err != nil {
		return
	}

	pss.Stream(cpeer, msg)
}

// Set the latency value by checking the UUID and peer match in the received 'pong' message.
func (pss *Server) optPong(cpeer peer.ID, uuidBytes []byte) {
	var newUUID uuid.NullUUID
	err := newUUID.UnmarshalBinary(uuidBytes)
	if err != nil {
		return
	}

	info := pss.apm.GetPeerInfo(cpeer)

	if info.uuid == newUUID.UUID {
		pss.apm.SetPeerInfoPongTime(cpeer, time.Now().UnixNano())
		info := pss.apm.GetPeerInfo(cpeer)

		// set latency value
		pss.apm.SetPeerInfoLatency(cpeer, time.Duration(info.pongTime-info.pingTime))

		// change uuid for next ping-pong
		pss.apm.SetPeerInfoUUID(cpeer, uuid.New())
	} else {
		// terminate the peer that engages in unexpected data exchange
		pss.host.Network().ClosePeer(cpeer)
	}
}

// Create a 'version' message and stream it to the peer that received the 'getversion' message.
func (pss *Server) optGetVersion(cpeer peer.ID) {
	msg, err := pss.omb.Version(pss.cfg.Version())
	if err != nil {
		return
	}

	pss.Stream(cpeer, msg)
}

// store peer version info in approved peers map
func (pss *Server) optVersion(cpeer peer.ID, bytes []byte) {
	pss.apm.SetPeerInfoVersion(cpeer, bytes)
}

// Create a 'peerlist' message and stream it to the peer that received the 'getpeerlist' message.
func (pss *Server) optGetPeerList(cpeer peer.ID) {
	msg, err := pss.omb.PeerList(pss.apm.ListApprovedPeerAddrs())
	if err != nil {
		return
	}

	pss.Stream(cpeer, msg)
}

// get peerlist from other peers
func (pss *Server) optPeerList(cpeer peer.ID, bytes []byte) {
	var addrs []peer.AddrInfo

	err := json.Unmarshal(bytes, &addrs)
	if err != nil {
		return
	}

	for _, addr := range addrs {
		if addr.ID == pss.self {
			continue
		}

		if !pss.apm.InPeers(addr.ID) {
			//TODO make test without mdns
			go pss.host.Connect(context.Background(), addr)
		}
	}
}

// process and transfer preconfirmation bids to the channel
func (pss *Server) optPreconfBid(cpeer peer.ID, bytes []byte) {
	pss.preconfCh <- bytes
}

// get self peer.ID
func (pss Server) ID() peer.ID { return pss.self }

// get approved peers on pubsub server
func (pss *Server) GetApprovedPeers() []peer.ID {
	return pss.apm.ListApprovedPeers()
}

// get approved gossip peers on pubsub server
func (pss *Server) GetApprovedGossipPeers() []peer.ID {
	return pss.apm.ListApprovedGossipPeers()
}

// when it listens to events, it applies necessary procedures based on incoming
// or outgoing connections
func (pss *Server) events(trackCh <-chan commons.ConnectionEvent) {
	var mutex = &sync.Mutex{}
	for event := range trackCh {

		// protect events
		mutex.Lock()
		eventCopy := event
		mutex.Unlock()

		switch eventCopy.Event {
		case commons.Connected:
			var retry = 5
			// create approve message and stream this message for new peer
			msg, err := pss.omb.Approve(pss.token)
			if err != nil {
				return
			}

			for i := 0; i < retry; i++ {
				err = pss.Stream(eventCopy.PeerID, msg)
				if err == nil {
					break
				}
			}

			go func() {
				checker := time.NewTicker(1000 * time.Millisecond)

				for i := 0; i < 10; i++ {
					// wait for it to join the verified peers
					<-checker.C

					if pss.apm.InPeers(eventCopy.PeerID) {
						// once the peer is connected,
						// send a message to get the version information
						msg, err = pss.omb.GetVersion()
						if err != nil {
							return
						}

						pss.Stream(eventCopy.PeerID, msg)
						break
					}
				}

				checker.Stop()

				if !pss.apm.InPeers(eventCopy.PeerID) {
					pss.host.Network().ClosePeer(eventCopy.PeerID)
					return
				}

				// generate a ping message with a unique UUID for each peer and append timestamps
				peerInfo := pss.apm.GetPeerInfo(eventCopy.PeerID)
				uuidBytes, err := peerInfo.uuid.MarshalBinary()
				if err != nil {
					return
				}

				msg, err = pss.omb.Ping(uuidBytes)
				if err != nil {
					return
				}

				pss.apm.SetPeerInfoPingTime(eventCopy.PeerID, time.Now().UnixNano())
				pss.Stream(eventCopy.PeerID, msg)

				// retrieve the peer list from the new node and expand the network connections
				msg, err = pss.omb.GetPeerList()
				if err != nil {
					return
				}

				pss.Stream(eventCopy.PeerID, msg)
			}()

		case commons.Disconnected:
			pss.apm.DelPeer(eventCopy.PeerID)
		}
	}
}

// periodically check the status of peers and update latency information
func (pss *Server) latencyUpdater() {

	// get latency check interval
	interval := pss.cfg.LatencyInterval()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		for peerID, info := range pss.apm.GetPeers() {
			uuidBytes, err := info.uuid.MarshalBinary()
			if err != nil {
				continue
			}

			msg, err := pss.omb.Ping(uuidBytes)
			if err != nil {
				continue
			}

			pss.apm.SetPeerInfoPingTime(peerID, time.Now().UnixNano())
			pss.Stream(peerID, msg)
		}
	}
}

// periodically update score information
func (pss *Server) scoreUpdater() {

	// sigmoid activation function
	sigmoid := func(x float64) float64 {
		return 1.0 / (1.0 + math.Exp(-x))
	}

	// get the score update interval from configuration
	interval := pss.cfg.ScoreInterval()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		// calculate and update score for each peer
		for peerID, info := range pss.apm.GetPeers() {
			// weights for each feature in the scoring formula
			w1, w2, w3 := 0.2, 0.3, 0.5

			start := info.getStart()
			stake := info.getStake()
			latency := info.getLatency()

			// normalize v1 by dividing milliseconds by the number of milliseconds in a day
			v1 := float64(time.Now().Sub(start).Milliseconds() / 864000000.)
			if v1 > 1 {
				v1 = 1.
			}

			// normalize v2 by dividing stake by 10^18
			v2 := float64(stake.Int64()) / float64(math.Pow(10, 18))
			if v2 > 1 {
				v2 = 1.
			}

			// normalize v3 by subtracting latency from 1000 milliseconds and taking the maximum of 0
			v3 := float64(1000 - latency.Milliseconds())
			if v3 < 0 {
				v3 = 0.
			}

			// calculate the score using the sigmoid activation function with the weighted values
			score := int(sigmoid(w1*v1+w2*v2+w3*v3) * 100.)

			// update the score for the peer
			pss.apm.SetPeerInfoScore(peerID, score)
		}
	}
}
