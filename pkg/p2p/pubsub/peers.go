package pubsub

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/peer"
)

type info struct {
	version string
	address common.Address
	stake   *big.Int
	//mode    commons.Mode
}

func (i *info) setVersion(version []byte) {
	i.version = string(version)
}

func (i *info) setAddress(address common.Address) {
	i.address = address
}

func (i *info) setStake(stake *big.Int) {
	i.stake = stake
}

//func (i *info) setMode(mode commons.Mode) {
//	i.mode = mode
//}

type approvedPeersMap struct {
	peers map[peer.ID]*info
	sync.RWMutex
}

func newApprovedPeersMap() *approvedPeersMap {
	return &approvedPeersMap{
		peers: make(map[peer.ID]*info, 1024),
	}
}

func (a *approvedPeersMap) AddPeer(peer peer.ID) {
	a.Lock()
	defer a.Unlock()
	a.peers[peer] = &info{}
}

func (a *approvedPeersMap) DelPeer(peer peer.ID) {
	a.Lock()
	defer a.Unlock()
	delete(a.peers, peer)
}

func (a approvedPeersMap) InPeers(peer peer.ID) bool {
	a.RLock()
	defer a.RUnlock()
	for k, _ := range a.peers {
		if k == peer {
			return true
		}
	}
	return false
}

//func (a approvedPeersMap) InSearchers(peer peer.ID) bool {
//	a.RLock()
//	defer a.RUnlock()
//	for k, v := range a.peers {
//		if k == peer && v.mode == commons.SEARCHER {
//			return true
//		}
//	}
//	return false
//}
//
//func (a approvedPeersMap) InBuilders(peer peer.ID) bool {
//	a.RLock()
//	defer a.RUnlock()
//	for k, v := range a.peers {
//		if k == peer && v.mode == commons.BUILDER {
//			return true
//		}
//	}
//	return false
//}

//func (a approvedPeersMap) GetApprovedSearchers() []peer.ID {
//	a.RLock()
//	defer a.RUnlock()
//
//	approvedPeers := []peer.ID{}
//	for k, v := range a.peers {
//		if v.mode == commons.SEARCHER {
//			approvedPeers = append(approvedPeers, k)
//		}
//	}
//
//	return approvedPeers
//}
//
//func (a approvedPeersMap) GetApprovedBuilders() []peer.ID {
//	a.RLock()
//	defer a.RUnlock()
//
//	approvedPeers := []peer.ID{}
//	for k, v := range a.peers {
//		if v.mode == commons.BUILDER {
//			approvedPeers = append(approvedPeers, k)
//		}
//	}
//
//	return approvedPeers
//}

func (a approvedPeersMap) GetPeers() map[peer.ID]*info {
	a.RLock()
	defer a.RUnlock()

	var peers = make(map[peer.ID]*info)
	for k, v := range a.peers {
		peers[k] = v
	}

	return peers
}

func (a approvedPeersMap) ListApprovedPeers() []peer.ID {
	a.RLock()
	defer a.RUnlock()

	approvedPeers := []peer.ID{}
	for k, _ := range a.peers {
		approvedPeers = append(approvedPeers, k)
	}

	return approvedPeers
}

// (version) set peer info options
func (a *approvedPeersMap) SetPeerInfoVersion(peer peer.ID, version []byte) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setVersion(version)
	}
}

// (address) set peer info options
func (a *approvedPeersMap) SetPeerInfoAddress(peer peer.ID, address common.Address) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setAddress(address)
	}
}

// (stake) set peer info options
func (a *approvedPeersMap) SetPeerInfoStake(peer peer.ID, stake *big.Int) {
	a.Lock()
	defer a.Unlock()
	if val, ok := a.peers[peer]; ok {
		val.setStake(stake)
	}
}

//// (mode) set peer info options
//func (a *approvedPeersMap) SetPeerInfoMode(peer peer.ID, mode commons.Mode) {
//	a.Lock()
//	defer a.Unlock()
//	a.peers[peer].setMode(mode)
//}
