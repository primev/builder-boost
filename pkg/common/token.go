package common

import (
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateToken(msg string, key *ecdsa.PrivateKey) (string, error) {
	fmt.Println("message is: ", msg)
	hm := crypto.Keccak256Hash([]byte(msg))
	signature, err := crypto.Sign(hm[:], key)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// Returns the authenticated address and the fail/success state of the token verification
func VerifyToken(token string, msg string) (common.Address, bool) {
	signature, err := base64.StdEncoding.DecodeString(token)
	fmt.Println("err1", err)
	if err != nil {
		return common.Address{}, false
	}
	fmt.Println("message is: ", msg)
	hm := crypto.Keccak256Hash([]byte(msg))
	pubkey, err := crypto.SigToPub(hm[:], signature)
	fmt.Println("err2", err)
	if err != nil {
		return common.Address{}, false
	}
	fmt.Println(crypto.PubkeyToAddress(*pubkey))
	// Remove recovery id
	// https://goethereumbook.org/signature-verify/
	return crypto.PubkeyToAddress(*pubkey), crypto.VerifySignature(crypto.FromECDSAPub(pubkey), hm[:], signature[:len(signature)-1])
}
