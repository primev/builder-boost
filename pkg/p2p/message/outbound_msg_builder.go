package message

import (
	"encoding/json"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

var _ OutboundMsgBuilder = (*outMsgBuilder)(nil)

type OutboundMsgBuilder interface {
	Approve([]byte) (OutboundMessage, error)
	Ping([]byte) (OutboundMessage, error)
	Pong([]byte) (OutboundMessage, error)
	GetVersion() (OutboundMessage, error)
	Version(string) (OutboundMessage, error)
	GetPeerList() (OutboundMessage, error)
	PeerList([]peer.AddrInfo) (OutboundMessage, error)
	PreconfBid([]byte) (OutboundMessage, error)
}

type outMsgBuilder struct {
	builder *msgBuilder
}

func NewOutboundBuilder() OutboundMsgBuilder {
	builder, err := newMsgBuilder()
	if err != nil {
		panic(err)
	}

	return &outMsgBuilder{
		builder: builder,
	}
}

func (b *outMsgBuilder) Approve(data []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Approve,
		time.Now().UnixNano(),
		data,
	)
}

func (b *outMsgBuilder) Ping(uuid []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Ping,
		time.Now().UnixNano(),
		uuid,
	)
}

func (b *outMsgBuilder) Pong(uuid []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Pong,
		time.Now().UnixNano(),
		uuid,
	)
}

func (b *outMsgBuilder) GetVersion() (OutboundMessage, error) {
	return b.builder.createOutbound(
		GetVersion,
		time.Now().UnixNano(),
		[]byte(""),
	)
}

func (b *outMsgBuilder) Version(version string) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Version,
		time.Now().UnixNano(),
		[]byte(version),
	)
}

func (b *outMsgBuilder) GetPeerList() (OutboundMessage, error) {
	return b.builder.createOutbound(
		GetPeerList,
		time.Now().UnixNano(),
		[]byte(""),
	)
}

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

func (b *outMsgBuilder) Signature(signature []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Signature,
		time.Now().UnixNano(),
		signature,
	)
}

func (b *outMsgBuilder) BlockKey(key []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		BlockKey,
		time.Now().UnixNano(),
		key,
	)
}

func (b *outMsgBuilder) Bundle(bundle []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Bundle,
		time.Now().UnixNano(),
		bundle,
	)
}

func (b *outMsgBuilder) PreconfBid(bid []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		PreconfBid,
		time.Now().UnixNano(),
		bid,
	)
}
