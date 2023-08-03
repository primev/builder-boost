package message

type Op byte

const (
	UNDEFINED Op = iota
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
	PreconfirmationBid
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
	case PreconfirmationBid:
		return "preconfirmationbid"
	default:
		return "Unknown Op"
	}
}
