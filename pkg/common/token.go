package common

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateToken(msg string, key *ecdsa.PrivateKey) (string, error) {
	hm := crypto.Keccak256Hash([]byte(msg))
	signature, err := crypto.Sign(hm[:], key)
	if err != nil {
		return "", err
	}

	return string(signature), nil
}

// Returns the authenticated address and the fail/success state of the token verification
func VerifyToken(token string, msg string) (common.Address, bool) {
	hm := crypto.Keccak256Hash([]byte(msg))
	pubkey, err := crypto.SigToPub(hm[:], []byte(token))
	if err != nil {
		return common.Address{}, false
	}

	verificationState := crypto.VerifySignature(crypto.CompressPubkey(pubkey), hm[:], []byte(token))

	return crypto.PubkeyToAddress(*pubkey), verificationState
}
