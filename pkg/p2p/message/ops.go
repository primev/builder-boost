package message

type Op byte

const (
	UNDEFINED Op = iota
	Authentication
	Ping
	Pong
	GetVersion
	Version
	GetPeerList
	PeerList
	PreconfirmationBid
)

func (op Op) String() string {
	switch op {
	case Authentication:
		return "authentication"
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
	case PreconfirmationBid:
		return "preconfirmationbid"
	default:
		return "Unknown Op"
	}
}
