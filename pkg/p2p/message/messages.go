package message

import (
	"encoding/json"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

//var (
//	InboundMessage  = &inboundMessage{}
//	OutboundMessage = &outboundMessage{}
//)

var (
	_ InboundMessage  = (*inboundMessage)(nil)
	_ OutboundMessage = (*outboundMessage)(nil)
)

type InboundMessage interface {
	Peer() peer.ID
	Op() Op
	Timestamp() int64
	Bytes() []byte
}

type inboundMessage struct {
	op        Op
	peer      peer.ID
	timestamp int64
	bytes     []byte
}

func (inMsg *inboundMessage) Peer() peer.ID {
	return inMsg.peer
}

func (inMsg *inboundMessage) Op() Op {
	return inMsg.op
}

func (inMsg *inboundMessage) Timestamp() int64 {
	return inMsg.timestamp
}

func (inMsg *inboundMessage) Bytes() []byte {
	return inMsg.bytes
}

func (inMsg *inboundMessage) SetPeer(peer peer.ID) {
	inMsg.peer = peer
}

func (inMsg *inboundMessage) UnmarshalJSON(b []byte) error {
	var tmp struct {
		Op        Op
		Peer      peer.ID
		Timestamp int64
		Bytes     []byte
	}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	inMsg.op = tmp.Op
	inMsg.peer = tmp.Peer
	inMsg.timestamp = tmp.Timestamp
	inMsg.bytes = tmp.Bytes
	return nil
}

func (inMsg *inboundMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Op        Op
		Timestamp int64
		Bytes     []byte
	}{
		Op:        inMsg.op,
		Timestamp: inMsg.timestamp,
		Bytes:     inMsg.bytes,
	})
}

type OutboundMessage interface {
	Op() Op
	Timestamp() int64
	Bytes() []byte
	SetOp(Op)
	SetTimestamp()
	SetBytes([]byte)
	UnmarshalJSON([]byte) error
	MarshalJSON() ([]byte, error)
}

type outboundMessage struct {
	op        Op     `json:"op"`
	timestamp int64  `json:"timestamp"`
	bytes     []byte `json:"bytes"`
}

func (outMsg *outboundMessage) Op() Op {
	return outMsg.op
}

func (outMsg *outboundMessage) Timestamp() int64 {
	return outMsg.timestamp
}

func (outMsg *outboundMessage) Bytes() []byte {
	return outMsg.bytes
}

func (outMsg *outboundMessage) SetOp(op Op) {
	outMsg.op = op
}

func (outMsg *outboundMessage) SetTimestamp() {
	outMsg.timestamp = time.Now().UnixNano()
}

func (outMsg *outboundMessage) SetBytes(bytes []byte) {
	outMsg.bytes = bytes
}

func (outMsg *outboundMessage) UnmarshalJSON(b []byte) error {
	var tmp struct {
		Op        Op
		Timestamp int64
		Bytes     []byte
	}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	outMsg.op = tmp.Op
	outMsg.timestamp = tmp.Timestamp
	outMsg.bytes = tmp.Bytes
	return nil
}

func (outMsg *outboundMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Op        Op
		Timestamp int64
		Bytes     []byte
	}{
		Op:        outMsg.op,
		Timestamp: outMsg.timestamp,
		Bytes:     outMsg.bytes,
	})
}
