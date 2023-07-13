package messages

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/peer"
)

type ApproveMsg struct {
	Address common.Address
	Peer    peer.ID
	Sig     []uint8
}

func (a *ApproveMsg) GetUnsignedMessage() []byte {
	return []byte(
		fmt.Sprintf(
			"%v:%v",
			a.Address.Hex(),
			a.Peer.String(),
		),
	)
}
