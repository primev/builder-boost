package node

import (
	"math/big"

	"github.com/libp2p/go-libp2p/core/control"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/lthibault/log"
	"github.com/multiformats/go-multiaddr"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/primev/builder-boost/pkg/p2p/commons"
	"github.com/primev/builder-boost/pkg/rollup"
)

type connectionAllowance int

const (
	Undecided connectionAllowance = iota
	DenyBannedPeer
	DenyNotEnoughStake
	Accept
)

var connectionAllowanceStrings = map[connectionAllowance]string{
	Undecided:          "Undecided",
	DenyBannedPeer:     "DenyBannedPeer",
	DenyNotEnoughStake: "DenyNotEnoughStake",
	Accept:             "Allow",
}

func (c connectionAllowance) isDeny() bool {
	return !(c == Accept || c == Undecided)
}

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
	blocker      *blocker
	logger       log.Logger
	minimalStake *big.Int
}

// newConnectionGater creates a new instance of ConnectionGater
func newConnectionGater(metrics *metrics, rollup rollup.Rollup, blocker *blocker, logger log.Logger, minimalStake *big.Int) ConnectionGater {
	return &connectionGater{
		metrics:      metrics,
		rollup:       rollup,
		blocker:      blocker,
		logger:       logger,
		minimalStake: minimalStake,
	}
}

// checkPeerTrusted determines the trust status of a peer
func (cg *connectionGater) checkPeerTrusted(p peer.ID) connectionAllowance {
	// TODO: Implement the logic to determine whether the peer is trusted or not
	return Undecided
}

// checkPeerBanned checks if a peer is banned and returns the appropriate connection allowance status
func (cg *connectionGater) checkPeerBanned(p peer.ID) connectionAllowance {

	// check if the peer is in the list of blocked peers, and deny the connection if found
	for _, peerID := range cg.blocker.list() {
		if p == peerID {
			return DenyBannedPeer
		}
	}

	// if the peer is not in the blocked list, allow the connection
	return Accept
}

// checkPeerStake checks if a peer has enough stake and returns the appropriate
// connection allowance status
func (cg *connectionGater) checkPeerStake(p peer.ID) connectionAllowance {
	pub, err := p.ExtractPublicKey()
	if err != nil {
		return DenyNotEnoughStake
	}

	// compressed sec256k1 public key bytes
	pubBytes, err := pub.Raw()
	if err != nil {
		return DenyNotEnoughStake
	}

	// extract eth address
	ethAddress := commons.Secp256k1CompressedBytesToEthAddress(pubBytes)

	// get stake
	stake, err := cg.rollup.GetMinimalStake(ethAddress)
	if err != nil {
		return DenyBannedPeer
	}

	// check minimal stake
	if stake.Cmp(cg.minimalStake) >= 0 {
		return Accept
	}

	// deny the connection if the stake is not enough.
	return DenyBannedPeer
}

// resolveConnectionAllowance resolves the connection allowance based on trusted and banned statuses
func resolveConnectionAllowance(
	trustedStatus connectionAllowance,
	bannedStatus connectionAllowance,
) connectionAllowance {
	// if the peer's trusted status is 'Undecided', resolve the connection allowance based on the banned status
	if trustedStatus == Undecided {
		return bannedStatus
	}
	return trustedStatus
}

// checks if a peer is allowed to dial/accept
func (cg *connectionGater) checkAllowedPeer(p peer.ID) connectionAllowance {
	return resolveConnectionAllowance(cg.checkPeerTrusted(p), cg.checkPeerBanned(p))
}

// InterceptPeerDial intercepts the process of dialing a peer
//
// all peer dialing attempts are allowed
func (cg *connectionGater) InterceptPeerDial(p peer.ID) bool {
	allowance := cg.checkAllowedPeer(p)
	if allowance.isDeny() {
		return false
	}

	return !cg.checkPeerStake(p).isDeny()
}

// InterceptAddrDial intercepts the process of dialing an address
//
// all address dialing attempts are allowed
// TODO rate limiter
func (cg *connectionGater) InterceptAddrDial(p peer.ID, addr ma.Multiaddr) bool {
	allowance := cg.checkAllowedPeer(p)
	if allowance.isDeny() {
		return false
	}

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
	allowance := cg.checkAllowedPeer(p)
	if allowance.isDeny() {
		return false
	}

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

// validateInboundConnection validates an inbound connection by extracting its
// public key and performing stake validation if the validation succeeds and
// the peer's stake is greater than minimal stake, the connection is allowed
// otherwise, the connection is rejected
func (cg *connectionGater) validateInboundConnection(p peer.ID, connMultiaddrs network.ConnMultiaddrs) bool {
	cg.metrics.IncomingConnectionCount.Inc()

	allowance := cg.checkPeerStake(p)
	if allowance.isDeny() {
		return false
	}

	return true
}

// validateOutboundConnection validates an outbound connection by extracting
// its public key and performing stake validation if the validation succeeds
// and the peer's stake is greater than minimal stake, the connection is
// allowed otherwise, the connection is rejected
func (cg *connectionGater) validateOutboundConnection(p peer.ID, connMultiaddrs network.ConnMultiaddrs) bool {
	cg.metrics.OutgoingConnectionCount.Inc()

	allowance := cg.checkPeerStake(p)
	if allowance.isDeny() {
		return false
	}

	return true
}
