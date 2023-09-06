package pubsub

import (
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/primev/builder-boost/pkg/p2p/commons"
)

type info struct {
	// peer type (builder or searcher)
	peerType commons.PeerType
	// join date of the peer
	joinDate time.Time
	// version information of the peer
	version string
	// address of the peer
	address common.Address
	// stake amount held by the peer
	stake *big.Int
	// addrs list
	hostAddrs []multiaddr.Multiaddr
	// peer uniq uuid
	uuid uuid.UUID
	// last ping time (out)
	pingTime int64
	// last pong time (in)
	pongTime int64
	// diff between ping and pong
	latency time.Duration
	// peer score
	score int
	// permission for gossip to peer
	gossip bool
	// protect info
	sync.RWMutex
}

// getMode returns the type of the peer.
func (i *info) getPeerType() commons.PeerType {
	i.RLock()
	defer i.RUnlock()

	return i.peerType
}

// getJoinDate returns the join date of the peer.
func (i *info) getJoinDate() time.Time {
	i.RLock()
	defer i.RUnlock()

	return i.joinDate
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

// getUUID returns the uuid of the peer.
func (i *info) getUUID() uuid.UUID {
	i.RLock()
	defer i.RUnlock()

	return i.uuid
}

// getPingTime returns the ping time of the peer.
func (i *info) getPingTime() int64 {
	i.RLock()
	defer i.RUnlock()

	return i.pingTime
}

// getPongTime returns the pong time of the peer.
func (i *info) getPongTime() int64 {
	i.RLock()
	defer i.RUnlock()

	return i.pongTime
}

// getLatency returns the latency of the peer.
func (i *info) getLatency() time.Duration {
	i.RLock()
	defer i.RUnlock()

	return i.latency
}

// getScore returns the score of the peer.
func (i *info) getScore() int {
	i.RLock()
	defer i.RUnlock()

	return i.score
}

// getGossip returns the gossip permisison of the peer.
func (i *info) getGossip() bool {
	i.RLock()
	defer i.RUnlock()

	return i.gossip
}

// setMode sets the mode of the peer.
func (i *info) setPeerType(peerType commons.PeerType) {
	i.Lock()
	defer i.Unlock()

	i.peerType = peerType
}

// setJoinDate sets the joinDate of the peer.
func (i *info) setJoinDate(joinDate time.Time) {
	i.Lock()
	defer i.Unlock()

	i.joinDate = joinDate
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

// setAddrs sets the host addrs of the peer.
func (i *info) setAddrs(addrs []multiaddr.Multiaddr) {
	i.Lock()
	defer i.Unlock()

	i.hostAddrs = addrs
}

// setUUID sets the uuid of the peer.
func (i *info) setUUID(uuid uuid.UUID) {
	i.Lock()
	defer i.Unlock()

	i.uuid = uuid
}

// setPingTime sets the pingTime of the peer.
func (i *info) setPingTime(pingTime int64) {
	i.Lock()
	defer i.Unlock()

	i.pingTime = pingTime
}

// setPongTime sets the pongTime of the peer.
func (i *info) setPongTime(pongTime int64) {
	i.Lock()
	defer i.Unlock()

	i.pongTime = pongTime
}

// setLatency sets the latency of the peer.
func (i *info) setLatency(latency time.Duration) {
	i.Lock()
	defer i.Unlock()

	i.latency = latency
}

// setScore sets the host addrs of the peer.
func (i *info) setScore(score int) {
	i.Lock()
	defer i.Unlock()

	i.score = score
}

// setGossip sets the gossip permission of the peer.
func (i *info) setGossip(gossip bool) {
	i.Lock()
	defer i.Unlock()

	i.gossip = gossip
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

// InBuilderPeers checks if a peer is in the approved builder peer.
func (a *approvedPeersMap) InBuilderPeers(peer peer.ID) bool {
	a.RLock()
	defer a.RUnlock()

	val, ok := a.peers[peer]
	if ok {
		return val.getPeerType() == commons.Builder
	}

	return false
}

// InSearcherPeers checks if a peer is in the approved searcher peer.
func (a *approvedPeersMap) InSearcherPeers(peer peer.ID) bool {
	a.RLock()
	defer a.RUnlock()

	val, ok := a.peers[peer]
	if ok {
		return val.getPeerType() == commons.Searcher
	}

	return false
}

// GetPeerInfo returns the information of a specific peer.
func (a *approvedPeersMap) GetPeerInfo(peer peer.ID) *info {
	a.RLock()
	defer a.RUnlock()

	if val, ok := a.peers[peer]; ok {
		infoCopy := &info{
			peerType:  val.getPeerType(),
			joinDate:  val.getJoinDate(),
			version:   val.getVersion(),
			address:   val.getAddress(),
			stake:     val.getStake(),
			hostAddrs: val.getAddrs(),
			uuid:      val.getUUID(),
			pingTime:  val.getPingTime(),
			pongTime:  val.getPongTime(),
			latency:   val.getLatency(),
			score:     val.getScore(),
			gossip:    val.getGossip(),
		}

		return infoCopy
	}

	return nil
}

// GetPeers returns a map of all the approved peers.
func (a *approvedPeersMap) GetPeers() map[peer.ID]*info {
	a.RLock()
	defer a.RUnlock()

	var peers = make(map[peer.ID]*info)
	for k, v := range a.peers {
		infoCopy := &info{
			peerType:  v.getPeerType(),
			joinDate:  v.getJoinDate(),
			version:   v.getVersion(),
			address:   v.getAddress(),
			stake:     v.getStake(),
			hostAddrs: v.getAddrs(),
			uuid:      v.getUUID(),
			pingTime:  v.getPingTime(),
			pongTime:  v.getPongTime(),
			latency:   v.getLatency(),
			score:     v.getScore(),
			gossip:    v.getGossip(),
		}

		peers[k] = infoCopy
	}

	return peers
}

// GetPeerIDs returns a list of all the approved peer IDs.
func (a *approvedPeersMap) GetPeerIDs() []peer.ID {
	a.RLock()
	defer a.RUnlock()

	peerIDs := []peer.ID{}
	for k := range a.peers {
		peerIDs = append(peerIDs, k)
	}

	return peerIDs
}

// GetBuilderPeerIDs returns a list of all the approved builder peer IDs.
func (a *approvedPeersMap) GetBuilderPeerIDs() []peer.ID {
	a.RLock()
	defer a.RUnlock()

	peerIDs := []peer.ID{}
	for k, val := range a.peers {
		if val.getPeerType() == commons.Builder {
			peerIDs = append(peerIDs, k)
		}
	}

	return peerIDs
}

// GetSearcherPeerIDs returns a list of all the approved searcher peer IDs.
func (a *approvedPeersMap) GetSearcherPeerIDs() []peer.ID {
	a.RLock()
	defer a.RUnlock()

	peerIDs := []peer.ID{}
	for k, val := range a.peers {
		if val.getPeerType() == commons.Searcher {
			peerIDs = append(peerIDs, k)
		}
	}

	return peerIDs
}

// GetGossipPeerIDs returns a list of all the approved gossip peer IDs.
func (a *approvedPeersMap) GetGossipPeerIDs() []peer.ID {
	a.RLock()
	defer a.RUnlock()

	peerIDs := []peer.ID{}
	for k, v := range a.peers {
		if v.getGossip() {
			peerIDs = append(peerIDs, k)
		}
	}

	return peerIDs
}

// GetGossipBuilderPeerIDs returns a list of all the approved builder gossip peer IDs.
func (a *approvedPeersMap) GetGossipBuilderPeerIDs() []peer.ID {
	a.RLock()
	defer a.RUnlock()

	peerIDs := []peer.ID{}
	for k, v := range a.peers {
		if v.getGossip() && (v.getPeerType() == commons.Builder) {
			peerIDs = append(peerIDs, k)
		}
	}

	return peerIDs
}

// GetGossipSearcherPeerIDs returns a list of all the approved searcher gossip peer IDs.
func (a *approvedPeersMap) GetGossipSearcherPeerIDs() []peer.ID {
	a.RLock()
	defer a.RUnlock()

	peerIDs := []peer.ID{}
	for k, v := range a.peers {
		if v.getGossip() && (v.getPeerType() == commons.Searcher) {
			peerIDs = append(peerIDs, k)
		}
	}

	return peerIDs
}

// GetBuilderPeerAddrs returns a list of all the approved peer addrs.
func (a *approvedPeersMap) GetBuilderPeerAddrs() []peer.AddrInfo {
	a.RLock()
	defer a.RUnlock()

	peerAddrs := []peer.AddrInfo{}

	for k, v := range a.peers {
		if v.getPeerType() == commons.Builder {
			addr := peer.AddrInfo{
				ID:    k,
				Addrs: v.getAddrs(),
			}
			peerAddrs = append(peerAddrs, addr)
		}

	}

	return peerAddrs
}

// GetGossipPeers returns a map of all the approved gossip peers.
func (a *approvedPeersMap) GetGossipPeers() map[peer.ID]*info {
	a.RLock()
	defer a.RUnlock()

	var peers = make(map[peer.ID]*info)
	for k, v := range a.peers {
		if v.getGossip() {
			infoCopy := &info{
				peerType:  v.getPeerType(),
				joinDate:  v.getJoinDate(),
				version:   v.getVersion(),
				address:   v.getAddress(),
				stake:     v.getStake(),
				hostAddrs: v.getAddrs(),
				uuid:      v.getUUID(),
				pingTime:  v.getPingTime(),
				pongTime:  v.getPongTime(),
				latency:   v.getLatency(),
				score:     v.getScore(),
				gossip:    v.getGossip(),
			}

			peers[k] = infoCopy
		}
	}

	return peers
}

// SetPeerInfoPeerType sets the type of a peer.
func (a *approvedPeersMap) SetPeerInfoPeerType(peer peer.ID, peerType commons.PeerType) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setPeerType(peerType)
	}
}

// SetPeerInfoJoinDate sets the join date of a peer.
func (a *approvedPeersMap) SetPeerInfoJoinDate(peer peer.ID, joinDate time.Time) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setJoinDate(joinDate)
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

// SetPeerInfoUUID sets the uuid of a peer.
func (a *approvedPeersMap) SetPeerInfoUUID(peer peer.ID, uuid uuid.UUID) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setUUID(uuid)
	}
}

// SetPeerInfoPingTime sets the ping time of a peer.
func (a *approvedPeersMap) SetPeerInfoPingTime(peer peer.ID, pingTime int64) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setPingTime(pingTime)
	}
}

// SetPeerInfoPongTime sets the pong time of a peer.
func (a *approvedPeersMap) SetPeerInfoPongTime(peer peer.ID, pongTime int64) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setPongTime(pongTime)
	}
}

// SetPeerInfoLatency sets the latency of a peer.
func (a *approvedPeersMap) SetPeerInfoLatency(peer peer.ID, latency time.Duration) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setLatency(latency)
	}
}

// SetPeerInfoScore sets the score of a peer.
func (a *approvedPeersMap) SetPeerInfoScore(peer peer.ID, score int) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setScore(score)
	}
}

// SetPeerInfoGossip sets the gossip permission of a peer.
func (a *approvedPeersMap) SetPeerInfoGossip(peer peer.ID, gossip bool) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setGossip(gossip)
	}
}
