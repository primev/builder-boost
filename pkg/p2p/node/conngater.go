package node

import (
	"math/big"

	"github.com/libp2p/go-libp2p/core/control"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/primev/builder-boost/pkg/p2p/commons"
	"github.com/primev/builder-boost/pkg/rollup"
)

type ConnectionGater interface {
	// InterceptPeerDial intercepts peer dialing
	InterceptPeerDial(p peer.ID) (allow bool)
	// InterceptAddrDial intercepts address dialing
	InterceptAddrDial(peer.ID, multiaddr.Multiaddr) (allow bool)
	// InterceptAccept intercepts connection acceptance
	InterceptAccept(network.ConnMultiaddrs) (allow bool)
	// InterceptSecured intercepts secured connection
	InterceptSecured(network.Direction, peer.ID, network.ConnMultiaddrs) (allow bool)
	// InterceptUpgraded intercepts upgraded connection
	InterceptUpgraded(network.Conn) (allow bool, reason control.DisconnectReason)
}

type connectionGater struct {
	metrics      *metrics
	rollup       rollup.Rollup
	minimalStake *big.Int
}

// newConnectionGater creates a new instance of ConnectionGater
func newConnectionGater(metrics *metrics, rollup rollup.Rollup, minimalStake *big.Int) ConnectionGater {
	return &connectionGater{
		metrics:      metrics,
		rollup:       rollup,
		minimalStake: minimalStake,
	}
}

// InterceptPeerDial intercepts the process of dialing a peer
//
// all peer dialing attempts are allowed
func (cg *connectionGater) InterceptPeerDial(p peer.ID) bool {
	return true
}

// InterceptAddrDial intercepts the process of dialing an address
//
// all address dialing attempts are allowed
// TODO rate limiter
func (cg *connectionGater) InterceptAddrDial(p peer.ID, addr ma.Multiaddr) bool {
	return true
}

// InterceptAccept intercepts the process of accepting a connection
//
// all connection acceptance attempts are allowed
func (cg *connectionGater) InterceptAccept(connMultiaddrs network.ConnMultiaddrs) bool {
	return true
}

// InterceptSecured intercepts a secured connection, regardless of its direction (inbound/outbound)
func (cg *connectionGater) InterceptSecured(dir network.Direction, p peer.ID, connMultiaddrs network.ConnMultiaddrs) bool {
	// note: we are indifferent to the direction (inbound/outbound)
	// if you want to manipulate (inbound/outbound) connections, make the change
	// note:if it is desired to not establish a connection with a peer,
	// ensure that it is rejected in both incoming and outgoing connections
	if dir == network.DirInbound {
		return cg.validateInboundConnection(p, connMultiaddrs)
	} else {
		return cg.validateOutboundConnection(p, connMultiaddrs)
	}

	return true
}

// InterceptUpgraded intercepts the process of upgrading a connection
//
// all connection upgrade attempts are allowed
func (cg *connectionGater) InterceptUpgraded(conn network.Conn) (bool, control.DisconnectReason) {
	return true, control.DisconnectReason(0)
}

// validateInboundConnection validates an inbound connection by extracting its public key and performing stake validation
// if the validation succeeds and the peer's stake is greater than zero, the connection is allowed
// otherwise, the connection is rejected
func (cg *connectionGater) validateInboundConnection(p peer.ID, connMultiaddrs network.ConnMultiaddrs) bool {

	cg.metrics.IncomingConnectionCount.Inc()

	pub, err := p.ExtractPublicKey()
	if err != nil {
		return false
	}

	// compressed sec256k1 public key bytes
	pubBytes, err := pub.Raw()
	if err != nil {
		return false
	}

	// extract eth address
	ethAddress := commons.Secp256k1CompressedBytesToEthAddress(pubBytes)

	// get stake
	stake, err := cg.rollup.GetMinimalStake(ethAddress)
	if err != nil {
		return false
	}

	// check minimal stake
	if stake.Cmp(cg.minimalStake) >= 0 {
		return true
	}

	return false
}

// validateOutboundConnection validates an outbound connection by extracting its public key and performing stake validation
// if the validation succeeds and the peer's stake is greater than zero, the connection is allowed
// otherwise, the connection is rejected
func (cg *connectionGater) validateOutboundConnection(p peer.ID, connMultiaddrs network.ConnMultiaddrs) bool {

	cg.metrics.OutgoingConnectionCount.Inc()

	pub, err := p.ExtractPublicKey()
	if err != nil {
		return false
	}

	// compressed sec256k1 public key bytes
	pubBytes, err := pub.Raw()
	if err != nil {
		return false
	}

	// extract eth address
	ethAddress := commons.Secp256k1CompressedBytesToEthAddress(pubBytes)

	// get stake
	stake, err := cg.rollup.GetMinimalStake(ethAddress)
	if err != nil {
		return false
	}

	// check minimal stake
	if stake.Cmp(cg.minimalStake) >= 0 {
		return true
	}

	return false
}
