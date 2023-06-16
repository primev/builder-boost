package pubsub

import (
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type info struct {
	// start time of the peer
	start time.Time
	// version information of the peer
	version string
	// address of the peer
	address common.Address
	// stake amount held by the peer
	stake *big.Int
	// addrs list
	hostAddrs []multiaddr.Multiaddr
	// protect info
	sync.RWMutex
}

// getStart returns the start time of the peer.
func (i *info) getStart() time.Time {
	i.RLock()
	defer i.RUnlock()

	return i.start
}

// getVersion returns the version information of the peer.
func (i *info) getVersion() string {
	i.RLock()
	defer i.RUnlock()

	return i.version
}

// getAddress returns the address of the peer.
func (i *info) getAddress() common.Address {
	i.RLock()
	defer i.RUnlock()

	return i.address
}

// getStake returns the stake amount held by the peer.
func (i *info) getStake() *big.Int {
	i.RLock()
	defer i.RUnlock()

	return i.stake
}

// getAddrs returns the host addrs of the peer.
func (i *info) getAddrs() []multiaddr.Multiaddr {
	i.RLock()
	defer i.RUnlock()

	return i.hostAddrs
}

// setStart sets the start time of the peer.
func (i *info) setStart(start time.Time) {
	i.Lock()
	defer i.Unlock()

	i.start = start
}

// setVersion sets the version information of the peer.
func (i *info) setVersion(version []byte) {
	i.Lock()
	defer i.Unlock()

	i.version = string(version)
}

// setAddress sets the address of the peer.
func (i *info) setAddress(address common.Address) {
	i.Lock()
	defer i.Unlock()

	i.address = address
}

// setStake sets the stake amount held by the peer.
func (i *info) setStake(stake *big.Int) {
	i.Lock()
	defer i.Unlock()

	i.stake = stake
}

// setAddrs returns the host addrs of the peer.
func (i *info) setAddrs(addrs []multiaddr.Multiaddr) {
	i.Lock()
	defer i.Unlock()

	i.hostAddrs = addrs
}

type approvedPeersMap struct {
	peers map[peer.ID]*info
	sync.RWMutex
}

// newApprovedPeersMap creates a new instance of the approvedPeersMap.
func newApprovedPeersMap() *approvedPeersMap {
	return &approvedPeersMap{
		peers: make(map[peer.ID]*info, 1024),
	}
}

// AddPeer adds a peer to the approved peers map.
func (a *approvedPeersMap) AddPeer(peer peer.ID) {
	a.Lock()
	defer a.Unlock()
	a.peers[peer] = &info{}
}

// DelPeer removes a peer from the approved peers map.
func (a *approvedPeersMap) DelPeer(peer peer.ID) {
	a.Lock()
	defer a.Unlock()
	delete(a.peers, peer)
}

// InPeers checks if a peer is in the approved peers map.
func (a *approvedPeersMap) InPeers(peer peer.ID) bool {
	a.RLock()
	defer a.RUnlock()

	_, ok := a.peers[peer]
	return ok
}

// GetPeers returns a map of all the approved peers.
func (a *approvedPeersMap) GetPeers() map[peer.ID]*info {
	a.RLock()
	defer a.RUnlock()

	var peers = make(map[peer.ID]*info)
	for k, v := range a.peers {
		infoCopy := &info{
			start:     v.getStart(),
			version:   v.getVersion(),
			address:   v.getAddress(),
			stake:     v.getStake(),
			hostAddrs: v.getAddrs(),
		}

		peers[k] = infoCopy
	}

	return peers
}

// ListApprovedPeers returns a list of all the approved peer IDs.
func (a *approvedPeersMap) ListApprovedPeers() []peer.ID {
	a.RLock()
	defer a.RUnlock()

	approvedPeers := []peer.ID{}
	for k := range a.peers {
		approvedPeers = append(approvedPeers, k)
	}

	return approvedPeers
}

// ListApprovedPeerAddrs returns a list of all the approved peer addrs.
func (a *approvedPeersMap) ListApprovedPeerAddrs() []peer.AddrInfo {
	a.RLock()
	defer a.RUnlock()

	approvedPeerAddrs := []peer.AddrInfo{}

	for k, v := range a.peers {
		addr := peer.AddrInfo{
			ID:    k,
			Addrs: v.getAddrs(),
		}

		approvedPeerAddrs = append(approvedPeerAddrs, addr)
	}

	return approvedPeerAddrs
}

// SetPeerInfoStart sets the start time of a peer.
func (a *approvedPeersMap) SetPeerInfoStart(peer peer.ID, start time.Time) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setStart(start)
	}
}

// SetPeerInfoVersion sets the version information of a peer.
func (a *approvedPeersMap) SetPeerInfoVersion(peer peer.ID, version []byte) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setVersion(version)
	}
}

// SetPeerInfoAddress sets the address of a peer.
func (a *approvedPeersMap) SetPeerInfoAddress(peer peer.ID, address common.Address) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setAddress(address)
	}
}

// SetPeerInfoStake sets the stake amount held by a peer.
func (a *approvedPeersMap) SetPeerInfoStake(peer peer.ID, stake *big.Int) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setStake(stake)
	}
}

// SetPeerInfoAddrs sets the host addrs of a peer.
func (a *approvedPeersMap) SetPeerInfoAddrs(peer peer.ID, addrs []multiaddr.Multiaddr) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setAddrs(addrs)
	}
}
