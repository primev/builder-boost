package node

import (
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/lthibault/log"
	"github.com/primev/builder-boost/pkg/p2p/commons"
)

// blocker struct manages blocked peers and their connections
type blocker struct {
	metrics      *metrics
	logger       log.Logger
	blockedPeers map[peer.ID]string
	mux          sync.RWMutex
}

// newBlocker creates a new blocker instance
func newBlocker(metrics *metrics, logger log.Logger) *blocker {
	return &blocker{
		metrics:      metrics,
		logger:       logger,
		blockedPeers: make(map[peer.ID]string),
	}
}

// make sure the blocked peer connection is closed
// add adds a peer to the blockedPeers map and increments metrics
func (b *blocker) add(peerID peer.ID, reason string) {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.blockedPeers[peerID] = reason
	b.metrics.BlockedPeerCount.Inc()
	b.metrics.CurrentBlockedPeerCount.Set(float64(len(b.blockedPeers)))
	b.metrics.BlockedPeers.WithLabelValues(peerID.String(), reason).Set(1)

	b.logger.With(log.F{
		"service":  "p2p blocker",
		"peer id":  peerID.String(),
		"reason":   reason,
		"log time": commons.GetNow(),
	}).Warn("peer blocked!")

}

// del removes a peer from the blockedPeers map and updates metrics.
func (b *blocker) del(peerID peer.ID) {
	b.mux.Lock()
	defer b.mux.Unlock()

	delete(b.blockedPeers, peerID)
	b.metrics.CurrentBlockedPeerCount.Set(float64(len(b.blockedPeers)))

	b.logger.With(log.F{
		"service":  "p2p blocker",
		"peer id":  peerID.String(),
		"log time": commons.GetNow(),
	}).Warn("peer block removed!")
}

// list returns a list of blocked peer IDs
func (b *blocker) list() []peer.ID {
	b.mux.RLock()
	defer b.mux.RUnlock()

	var peers []peer.ID
	for peerID, _ := range b.blockedPeers {
		peers = append(peers, peerID)
	}
	return peers
}
