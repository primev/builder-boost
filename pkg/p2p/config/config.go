package config

import (
	"time"

	"github.com/multiformats/go-multiaddr"
)

type Config struct {
	// p2p last version
	version string
	// DiscoveryInterval is how often we re-publish our mDNS records.
	discoveryInterval time.Duration
	// DiscoveryServiceTag is used in our mDNS advertisements to discover other peers.
	discoveryServiceTag string
	// Topic name for DHT
	topic string
	// General peer communication topic
	// @iowar: If there arises a need for a new topic, a name change should be made
	pubSubTopic string
	// default stream protocol
	peerStreamProto string
	// bootstrap peers
	bootstrapPeers []multiaddr.Multiaddr
}

type ConfigOption func(*Config)

// New creates a new Config with default values and applies the provided options
func New(options ...ConfigOption) *Config {
	cfg := &Config{
		version:             "0.0.1",
		discoveryInterval:   time.Hour,
		discoveryServiceTag: "PRIMEV-0.0.1",
		topic:               "PRIMEV-TEST-45",
		pubSubTopic:         "PRIMEVTOPIC",
		peerStreamProto:     "/primev/stream-0.1",
	}

	for _, option := range options {
		option(cfg)
	}

	return cfg
}

// WithVersion sets the version option for Config
func WithVersion(version string) ConfigOption {
	return func(cfg *Config) {
		cfg.version = version
	}
}

// WithDiscoveryInterval sets the discovery interval option for Config
func WithDiscoveryInterval(interval time.Duration) ConfigOption {
	return func(cfg *Config) {
		cfg.discoveryInterval = interval
	}
}

// WithDiscoveryServiceTag sets the discovery service tag option for Config
func WithDiscoveryServiceTag(serviceTag string) ConfigOption {
	return func(cfg *Config) {
		cfg.discoveryServiceTag = serviceTag
	}
}

// WithTopic sets the topic option for Config
func WithTopic(topic string) ConfigOption {
	return func(cfg *Config) {
		cfg.topic = topic
	}
}

// WithPubSubTopic sets the PubSub topic option for Config
func WithPubSubTopic(pubSubTopic string) ConfigOption {
	return func(cfg *Config) {
		cfg.pubSubTopic = pubSubTopic
	}
}

// WithPeerStreamProto sets the peer stream protocol option for Config
func WithPeerStreamProto(proto string) ConfigOption {
	return func(cfg *Config) {
		cfg.peerStreamProto = proto
	}
}

// WithBootstrapPeers is used to add custom bootstrap peers to the Config configuration
func WithBootstrapPeers(peers []multiaddr.Multiaddr) ConfigOption {
	return func(cfg *Config) {
		cfg.bootstrapPeers = peers
	}
}

// Version returns the version from Config
func (cfg *Config) Version() string {
	return cfg.version
}

// DiscoveryInterval returns the discovery interval from Config
func (cfg *Config) DiscoveryInterval() time.Duration {
	return cfg.discoveryInterval
}

// DiscoveryServiceTag returns the discovery service tag from Config
func (cfg *Config) DiscoveryServiceTag() string {
	return cfg.discoveryServiceTag
}

// Topic returns the topic from Config
func (cfg *Config) Topic() string {
	return cfg.topic
}

// PubSubTopic returns the PubSub topic from Config
func (cfg *Config) PubSubTopic() string {
	return cfg.pubSubTopic
}

// PeerStreamProto returns the peer stream protocol from Config
func (cfg *Config) PeerStreamProto() string {
	return cfg.peerStreamProto
}

// BootstrapPeers returns the bootstrap peers from Config
func (cfg *Config) BootstrapPeers() []multiaddr.Multiaddr {
	return cfg.bootstrapPeers
}