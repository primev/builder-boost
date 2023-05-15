package rollup

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type State struct {
	LatestProcessedBlock uint64                                       `json:"latestProcessedBlock"`
	Stakes               map[common.Address]map[common.Address]BigInt `json:"stakes"`
	MinimalStakes        map[common.Address]BigInt                    `json:"minimalStakes"`
}

type BigInt struct {
	big.Int
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *BigInt) UnmarshalJSON(p []byte) error {
	if string(p) == "null" {
		return nil
	}
	var z big.Int
	_, ok := z.SetString(string(p), 10)
	if !ok {
		return fmt.Errorf("not a valid big integer: %s", p)
	}
	b.Int = z
	return nil
}
