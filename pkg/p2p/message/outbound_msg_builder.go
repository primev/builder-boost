package message

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/primev/builder-boost/pkg/p2p/commons"
)

var _ OutboundMsgBuilder = (*outMsgBuilder)(nil)

type OutboundMsgBuilder interface {
	Authentication([]byte) (OutboundMessage, error)
	Ping() (OutboundMessage, error)
	Pong() (OutboundMessage, error)
	GetVersion() (OutboundMessage, error)
	Version(string) (OutboundMessage, error)
	GetPeerList() (OutboundMessage, error)
	PeerList([]peer.ID) (OutboundMessage, error)
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

func (b *outMsgBuilder) Ping() (OutboundMessage, error) {
	return b.builder.createOutbound(
		Ping,
		time.Now().UnixNano(),
		[]byte(""),
	)
}

func (b *outMsgBuilder) Pong() (OutboundMessage, error) {
	return b.builder.createOutbound(
		Pong,
		time.Now().UnixNano(),
		[]byte(""),
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

func (b *outMsgBuilder) PeerList(peers []peer.ID) (OutboundMessage, error) {
	return b.builder.createOutbound(
		PeerList,
		time.Now().UnixNano(),
		commons.PeersToBytes(peers),
	)
}
