package message

import (
	"encoding/json"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

var _ OutboundMsgBuilder = (*outMsgBuilder)(nil)

type OutboundMsgBuilder interface {
	Authentication([]byte) (OutboundMessage, error)
	Ping([]byte) (OutboundMessage, error)
	Pong([]byte) (OutboundMessage, error)
	GetVersion() (OutboundMessage, error)
	Version(string) (OutboundMessage, error)
	GetPeerList() (OutboundMessage, error)
	PeerList([]peer.AddrInfo) (OutboundMessage, error)
	PreconfirmationBid([]byte) (OutboundMessage, error)
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

func (b *outMsgBuilder) Authentication(data []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		Authentication,
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

func (b *outMsgBuilder) PreconfirmationBid(bid []byte) (OutboundMessage, error) {
	return b.builder.createOutbound(
		PreconfirmationBid,
		time.Now().UnixNano(),
		bid,
	)
}
