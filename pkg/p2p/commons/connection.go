package commons

import "github.com/libp2p/go-libp2p/core/peer"

type Event int

const (
	UndefinedE Event = iota
	Connected
	Disconnected
)

type ConnectionEvent struct {
	PeerID   peer.ID
	PeerType PeerType
	Event    Event
}
