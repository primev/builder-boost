package discover

import (
	"context"

	"github.com/lthibault/log"
	"github.com/primev/builder-boost/pkg/p2p/commons"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

// logging only mdns and dht discover
// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host

	ctx context.Context
	log log.Logger
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(p peer.AddrInfo) {
	n.log.With(log.F{
		"caller":  commons.GetCallerName(),
		"date":    commons.GetNow(),
		"service": "p2p notify",
		"peer":    p.ID.Pretty(),
	}).Info("discovered new peer")

	err := n.h.Connect(context.Background(), p)
	if err != nil {
		n.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p notify",
			"peer":    p.ID.Pretty(),
		}).Error(err)
	}
}
