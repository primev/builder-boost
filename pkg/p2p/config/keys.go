package config

import "time"

const (
	// p2p last version
	Version = "0.0.1"

	// DiscoveryInterval is how often we re-publish our mDNS records.
	DiscoveryInterval = time.Hour

	// DiscoveryServiceTag is used in our mDNS advertisements to discover other peers.
	DiscoveryServiceTag = "PRIMEV-0.0.1"

	// Topic name for DHT
	Topic = "PRIMEV-TEST-45"

	// General peer communication topic
	// @iowar: If there arises a need for a new topic, a name change should be made
	PubSubTopic = "PRIMEVTOPIC"

	// // PubSub topic b<>s s<>b s<>s b<>b
	// AllPubSubTopic = "PRIMEVTOPIC_ALL"
	//
	// // PubSub topic b<>b
	// B2bPubSubTopic = "PRIMEVTOPIC_B2B"
)

// default stream protocols
const (
	PeerStreamProto = "/primev/stream-0.1"
)
