package message

import (
	"encoding/json"

	"github.com/libp2p/go-libp2p/core/peer"
)

type msgBuilder struct{}

func newMsgBuilder() (*msgBuilder, error) {
	mb := &msgBuilder{}
	return mb, nil
}

func (mb *msgBuilder) createOutbound(op Op, timestamp int64, bytes []byte) (*outboundMessage, error) {
	return &outboundMessage{
		op:        op,
		timestamp: timestamp,
		bytes:     bytes,
	}, nil
}

func (mb *msgBuilder) parseInbound(peer peer.ID, bytes []byte) (*inboundMessage, error) {
	msg := new(inboundMessage)
	err := json.Unmarshal(bytes, msg)
	if err != nil {
		return nil, err
	}

	msg.peer = peer
	return msg, nil

}
