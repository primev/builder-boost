package config

import (
	"math/big"
	"time"

	"github.com/multiformats/go-multiaddr"
)

type Config struct {
	// p2p last version
	version string
	// p2p port
	p2pPort int
	// DiscoveryInterval is how often we re-publish our mDNS records.
	discoveryInterval time.Duration
	// DiscoveryServiceTag is used in our mDNS advertisements to discover other peers.
	discoveryServiceTag string
	// Topic name for DHT
	topic string
	// General peer communication topic
	// @iowar: If there arises a need for a new topic, a name change should be made
	// builder<>builder
	pubSubTopicB2B string
	// builder<>searcher
	pubSubTopicB2S string
	// default stream protocol b<>b
	peerStreamB2B string
	// default stream protocol b<>s
	peerStreamB2S string
	// bootstrap peers
	bootstrapPeers []multiaddr.Multiaddr
	// latency update interval
	latencyInterval time.Duration
	// score update interval
	scoreInterval time.Duration
	// builder minimal stake
	minimalStake *big.Int
	// metrics namespace
	metricsNamespace string
	// prometheus p2p metrics port
	metricsPort int
	// prometheus p2p metrics route
	metricsRoute string
}

type ConfigOption func(*Config)

// New creates a new Config with default values and applies the provided options
func New(options ...ConfigOption) *Config {
	cfg := &Config{
		version:             "0.0.1",
		p2pPort:             0,
		discoveryInterval:   time.Hour,
		discoveryServiceTag: "PRIMEV-0.0.1",
		topic:               "PRIMEV-TEST-45",
		pubSubTopicB2B:      "PRIMEVTOPICB2B",
		pubSubTopicB2S:      "PRIMEVTOPICB2S",
		peerStreamB2B:       "/primev/stream_b2b-0.1",
		peerStreamB2S:       "/primev/stream_b2s-0.1",
		latencyInterval:     time.Hour * 2,
		scoreInterval:       time.Hour * 3,
		minimalStake:        big.NewInt(0),
		metricsNamespace:    "builder_boost",
		metricsPort:         8081,
		metricsRoute:        "/metrics_p2p",
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

// WithP2PPort sets the version option for Config
func WithP2PPort(port int) ConfigOption {
	return func(cfg *Config) {
		cfg.p2pPort = port
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

// WithPubSubTopicB2B sets the PubSub topic option for Config
func WithPubSubTopicB2B(pubSubTopicB2B string) ConfigOption {
	return func(cfg *Config) {
		cfg.pubSubTopicB2B = pubSubTopicB2B
	}
}

// WithPubSubTopicB2S sets the PubSub topic option for Config
func WithPubSubTopicB2S(pubSubTopicB2S string) ConfigOption {
	return func(cfg *Config) {
		cfg.pubSubTopicB2S = pubSubTopicB2S
	}
}

// WithPeerStreamB2B sets the peer stream protocol option for Config
func WithPeerStreamB2B(proto string) ConfigOption {
	return func(cfg *Config) {
		cfg.peerStreamB2B = proto
	}
}

// WithPeerStreamB2S sets the peer stream protocol option for Config
func WithPeerStreamB2S(proto string) ConfigOption {
	return func(cfg *Config) {
		cfg.peerStreamB2S = proto
	}
}

// WithBootstrapPeers is used to add custom bootstrap peers to the Config configuration
func WithBootstrapPeers(peers []multiaddr.Multiaddr) ConfigOption {
	return func(cfg *Config) {
		cfg.bootstrapPeers = peers
	}
}

// WithLatencyInterval sets the latency interval option for Config
func WithLatencyInterval(interval time.Duration) ConfigOption {
	return func(cfg *Config) {
		cfg.latencyInterval = interval
	}
}

// WithScoreInterval sets the score interval option for Config
func WithScoreInterval(interval time.Duration) ConfigOption {
	return func(cfg *Config) {
		cfg.scoreInterval = interval
	}
}

// WithMinimalStake sets the minimal stake option for Config
func WithMinimalStake(minimalStake *big.Int) ConfigOption {
	return func(cfg *Config) {
		cfg.minimalStake = minimalStake
	}
}

// WithMetricsNamespace sets the metrics namespace option for Config
func WithMetricsNamespace(metricsNamespace string) ConfigOption {
	return func(cfg *Config) {
		cfg.metricsNamespace = metricsNamespace
	}
}

// WithMetricsPort sets the metrics port option for Config
func WithMetricsPort(metricsPort int) ConfigOption {
	return func(cfg *Config) {
		cfg.metricsPort = metricsPort
	}
}

// WithMetricsRoute sets the metrics route option for Config
func WithMetricsRoute(metricsRoute string) ConfigOption {
	return func(cfg *Config) {
		cfg.metricsRoute = metricsRoute
	}
}

// Version returns the version from Config
func (cfg *Config) Version() string {
	return cfg.version
}

// P2PPort returns the p2p port from Config
func (cfg *Config) P2PPort() int {
	return cfg.p2pPort
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

// PubSubTopicB2B returns the PubSub topic for b<>b from Config
func (cfg *Config) PubSubTopicB2B() string {
	return cfg.pubSubTopicB2B
}

// PubSubTopicB2S returns the PubSub topic for b<>s from Config
func (cfg *Config) PubSubTopicB2S() string {
	return cfg.pubSubTopicB2S
}

// PeerStreamB2B returns the peer stream protocol from Config
func (cfg *Config) PeerStreamB2B() string {
	return cfg.peerStreamB2B
}

// PeerStreamB2S returns the peer stream protocol from Config
func (cfg *Config) PeerStreamB2S() string {
	return cfg.peerStreamB2S
}

// BootstrapPeers returns the bootstrap peers from Config
func (cfg *Config) BootstrapPeers() []multiaddr.Multiaddr {
	return cfg.bootstrapPeers
}

// LatencyInterval returns the latency interval from Config
func (cfg *Config) LatencyInterval() time.Duration {
	return cfg.latencyInterval
}

// ScoreInterval returns the score interval from Config
func (cfg *Config) ScoreInterval() time.Duration {
	return cfg.scoreInterval
}

// MinimalStake returns the minimal stake from Config
func (cfg *Config) MinimalStake() *big.Int {
	return cfg.minimalStake
}

// MetricsNamespace returns the metrics namespace from Config
func (cfg *Config) MetricsNamespace() string {
	return cfg.metricsNamespace
}

// MetricsPort returns the metrics port from Config
func (cfg *Config) MetricsPort() int {
	return cfg.metricsPort
}

// MetricsRoute returns the metrics route from Config
func (cfg *Config) MetricsRoute() string {
	return cfg.metricsRoute
}
