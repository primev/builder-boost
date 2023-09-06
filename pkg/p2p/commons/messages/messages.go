package messages

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/primev/builder-boost/pkg/p2p/commons"
)

type ApproveMsg struct {
	PeerType commons.PeerType
	Address  common.Address
	Peer     peer.ID
	Sig      []uint8
}

func (a *ApproveMsg) GetUnsignedMessage() []byte {
	return []byte(
		fmt.Sprintf(
			"%v:%v:%v",
			a.PeerType,
			a.Address.Hex(),
			a.Peer.String(),
		),
	)
}

// P2P standard message type
type PeerMsg struct {
	Peer  peer.ID
	Bytes []byte
}
