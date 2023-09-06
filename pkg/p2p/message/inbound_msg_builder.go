package message

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

var _ InboundMsgBuilder = (*inMsgBuilder)(nil)

type InboundMsgBuilder interface {
	Parse(
		peer peer.ID,
		bytes []byte,
	) (InboundMessage, error)
}

type inMsgBuilder struct {
	builder *msgBuilder
}

func NewInboundBuilder() InboundMsgBuilder {
	builder, err := newMsgBuilder()
	if err != nil {
		panic(err)
	}

	return &inMsgBuilder{
		builder: builder,
	}
}

func (b *inMsgBuilder) Parse(peer peer.ID, bytes []byte) (InboundMessage, error) {
	return b.builder.parseInbound(peer, bytes)
}
