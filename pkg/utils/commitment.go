package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// GetCommitment calculates commitment by pseudo searcher and builder addresses
func GetCommitment(commitmentAddress, builder common.Address) common.Hash {
	return crypto.Keccak256Hash(commitmentAddress.Bytes(), builder.Bytes())
}
