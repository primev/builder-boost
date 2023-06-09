package node

import (
	"github.com/libp2p/go-libp2p/core/control"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/primev/builder-boost/pkg/rollup"
)

type ConnectionGater interface {
	InterceptPeerDial(p peer.ID) (allow bool)
	InterceptAddrDial(peer.ID, multiaddr.Multiaddr) (allow bool)
	InterceptAccept(network.ConnMultiaddrs) (allow bool)
	InterceptSecured(network.Direction, peer.ID, network.ConnMultiaddrs) (allow bool)
	InterceptUpgraded(network.Conn) (allow bool, reason control.DisconnectReason)
}

type connectionGater struct {
	rollup rollup.Rollup
}

func newConnectionGater(rollup rollup.Rollup) ConnectionGater {
	return &connectionGater{
		rollup: rollup,
	}
}

func (cg *connectionGater) InterceptPeerDial(p peer.ID) bool {
	return true
}

func (cg *connectionGater) InterceptAddrDial(p peer.ID, addr ma.Multiaddr) bool {
	return true
}

func (cg *connectionGater) InterceptAccept(connMultiaddrs network.ConnMultiaddrs) bool {
	return true
}

func (cg *connectionGater) InterceptSecured(dir network.Direction, p peer.ID, connMultiaddrs network.ConnMultiaddrs) bool {
	if dir == network.DirInbound {
		if validateInboundConnection(p, connMultiaddrs) {
			return true
		}
	} else {
		if validateOutgoingConnection(p, connMultiaddrs) {
			return true
		}
	}

	return true
}

func (cg *connectionGater) InterceptUpgraded(conn network.Conn) (bool, control.DisconnectReason) {
	return true, control.DisconnectReason(0)
}

func validateInboundConnection(p peer.ID, connMultiaddrs network.ConnMultiaddrs) bool {
	return true
}

func validateOutgoingConnection(p peer.ID, connMultiaddrs network.ConnMultiaddrs) bool {
	return true
}
