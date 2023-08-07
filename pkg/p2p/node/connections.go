package node

import (
	"sync"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/lthibault/log"
	"github.com/primev/builder-boost/pkg/p2p/commons"
)

// the purpose of the connection tracker packet is to trigger incoming and outgoing connections
// it was not designed for any other use case
type connectionTracker struct {
	metrics        *metrics
	log            log.Logger
	connectedPeers map[peer.ID]bool
	trackCh        chan commons.ConnectionEvent
	mux            sync.Mutex
}

func newConnectionTracker(metrics *metrics, log log.Logger) *connectionTracker {
	return &connectionTracker{
		metrics:        metrics,
		log:            log,
		connectedPeers: make(map[peer.ID]bool),
		trackCh:        make(chan commons.ConnectionEvent),
	}
}

func (ct *connectionTracker) handleConnected(net network.Network, conn network.Conn) {
	remotePeerID := conn.RemotePeer()

	ct.mux.Lock()
	defer ct.mux.Unlock()

	// Skip if the peer is already connected
	if ct.connectedPeers[remotePeerID] {
		return
	}

	ct.connectedPeers[remotePeerID] = true
	ct.sendConnected(remotePeerID)

	ct.log.With(log.F{
		"service":  "connection tracker",
		"peer":     remotePeerID.Pretty(),
		"log time": commons.GetNow(),
	}).Info("a new peer is connected")
}

func (ct *connectionTracker) handleDisconnected(net network.Network, conn network.Conn) {
	remotePeerID := conn.RemotePeer()

	ct.mux.Lock()
	defer ct.mux.Unlock()

	// Skip if the peer is already disconnected
	if !ct.connectedPeers[remotePeerID] {
		return
	}

	if net.Connectedness(remotePeerID) == network.Connected {
		return
	}

	delete(ct.connectedPeers, remotePeerID)
	ct.sendDisconnected(remotePeerID)

	ct.log.With(log.F{
		"service":  "connection tracker",
		"peer":     remotePeerID.Pretty(),
		"log time": commons.GetNow(),
	}).Info("a peer connection is disconnected")
}

func (ct *connectionTracker) sendConnected(peerID peer.ID) {
	// send the connected peer's information to the channel
	ct.trackCh <- commons.ConnectionEvent{
		PeerID: peerID,
		Event:  commons.Connected,
	}

	ct.metrics.ConnectedPeerCount.Inc()
	return
}

func (ct *connectionTracker) sendDisconnected(peerID peer.ID) {
	// send the disconnected peer's information to the channel
	ct.trackCh <- commons.ConnectionEvent{
		PeerID: peerID,
		Event:  commons.Disconnected,
	}

	ct.metrics.DisconnectedPeerCount.Inc()
	return
}
