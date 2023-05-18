package utils

import (
	"crypto/ecdsa"
	"crypto/x509"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// GetCommitment calculates commitment hash by builder private key and searcher address
func GetCommitment(key *ecdsa.PrivateKey, searcher common.Address) common.Hash {
	keyBytes, _ := x509.MarshalECPrivateKey(key)
	return crypto.Keccak256Hash(keyBytes, searcher.Bytes())
}
