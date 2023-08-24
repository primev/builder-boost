package message

import (
	"encoding/json"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

var _ OutboundMsgBuilder = (*outMsgBuilder)(nil)

// OutboundMsgBuilder defines an interface for constructing various outbound communication messages.
type OutboundMsgBuilder interface {
	// builder and searchers use this msg types
	Approve([]byte) (OutboundMessage, error)
	Ping([]byte) (OutboundMessage, error)
	Pong([]byte) (OutboundMessage, error)
	GetVersion() (OutboundMessage, error)
	Version(string) (OutboundMessage, error)
	GetPeerList() (OutboundMessage, error)
	PeerList([]peer.AddrInfo) (OutboundMessage, error)
	// only builders use this msg types
	Signature([]byte) (OutboundMessage, error)
	BlockKey([]byte) (OutboundMessage, error)
	Bundle([]byte) (OutboundMessage, error)
	PreconfBid([]byte) (OutboundMessage, error)
	// designed for searchers but builders use this msg types for provide communincation with searchers
	Bid([]byte) (OutboundMessage, error)
	Commitment([]byte) (OutboundMessage, error)
}

// outMsgBuilder defines a structure used for creating outbound communication messages.
type outMsgBuilder struct {
	builder *msgBuilder
}

// NewOutboundBuilder returns a new instance of OutboundMsgBuilder.
func NewOutboundBuilder() OutboundMsgBuilder {
	builder, err := newMsgBuilder()
	if err != nil {
		panic(err)
	}

	return &outMsgBuilder{
		builder: builder,
	}
}

// Approve creates an "Approve" message.
func (b *outMsgBuilder) Approve(data []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Approve,
		time.Now().UnixNano(),
		data,
	)
}

// Ping creates a "Ping" message.
func (b *outMsgBuilder) Ping(uuid []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Ping,
		time.Now().UnixNano(),
		uuid,
	)
}

// Pong creates a "Pong" message.
func (b *outMsgBuilder) Pong(uuid []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Pong,
		time.Now().UnixNano(),
		uuid,
	)
}

// GetVersion creates a "GetVersion" message.
func (b *outMsgBuilder) GetVersion() (OutboundMessage, error) {
	return b.builder.createOutbound(
		GetVersion,
		time.Now().UnixNano(),
		[]byte(""),
	)
}

// Version creates a "Version" message with the specified version information.
func (b *outMsgBuilder) Version(version string) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Version,
		time.Now().UnixNano(),
		[]byte(version),
	)
}

// GetPeerList creates a "GetPeerList" message.
func (b *outMsgBuilder) GetPeerList() (OutboundMessage, error) {
	return b.builder.createOutbound(
		GetPeerList,
		time.Now().UnixNano(),
		[]byte(""),
	)
}

// PeerList creates a "PeerList" message with the given address information.
func (b *outMsgBuilder) PeerList(addrs []peer.AddrInfo) (OutboundMessage, error) {
	data, err := json.Marshal(&addrs)
	if err != nil {
		return nil, err
	}

	return b.builder.createOutbound(
		PeerList,
		time.Now().UnixNano(),
		data,
	)
}

// Signature creates a "Signature" message with the specified signature data.
func (b *outMsgBuilder) Signature(signature []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Signature,
		time.Now().UnixNano(),
		signature,
	)
}

// BlockKey creates a "BlockKey" message with the specified key data.
func (b *outMsgBuilder) BlockKey(key []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		BlockKey,
		time.Now().UnixNano(),
		key,
	)
}

// Bundle creates a "Bundle" message with the specified bundle data.
func (b *outMsgBuilder) Bundle(bundle []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Bundle,
		time.Now().UnixNano(),
		bundle,
	)
}

// PreconfBid creates a "PreconfBid" message with the specified bid data.
func (b *outMsgBuilder) PreconfBid(bid []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		PreconfBid,
		time.Now().UnixNano(),
		bid,
	)
}

// Bid creates a "Bid" message with the specified bid data.
func (b *outMsgBuilder) Bid(bid []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Bid,
		time.Now().UnixNano(),
		bid,
	)
}

// Commitment creates a "commitment" message with the specified commitment data.
func (b *outMsgBuilder) Commitment(commitment []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Commitment,
		time.Now().UnixNano(),
		commitment,
	)
}
