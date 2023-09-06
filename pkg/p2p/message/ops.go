package message

type Op byte

const (
	UNDEFINED_op Op = iota
	Approve
	Ping
	Pong
	GetVersion
	Version
	GetPeerList
	PeerList
	Signature
	BlockKey
	Bundle
	PreconfBid
	Bid
	Commitment
)

func (op Op) String() string {
	switch op {
	case Approve:
		return "approve"
	case Ping:
		return "ping"
	case Pong:
		return "pong"
	case GetVersion:
		return "getversion"
	case Version:
		return "version"
	case GetPeerList:
		return "getpeerlist"
	case PeerList:
		return "peerlist"
	case Signature:
		return "signature"
	case BlockKey:
		return "blockkey"
	case Bundle:
		return "bundle"
	case PreconfBid:
		return "preconfbid"
	case Bid:
		return "bid"
	case Commitment:
		return "commitment"

	default:
		return "Unknown Op"
	}
}
