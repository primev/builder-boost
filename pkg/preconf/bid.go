package preconf

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

var (
	ErrAlreadySignedBid     = errors.New("already contains hash or signature")
	ErrMissingHashSignature = errors.New("missing hash or signature")
	ErrInvalidSignature     = errors.New("signature is not valid")
	ErrInvalidHash          = errors.New("bidhash doesn't match bid payload")
)

type UnsignedPreConfBid struct {
	TxnHash string   `json:"txnHash"`
	Bid     *big.Int `json:"bid"`
	// UUID    string `json:"uuid"` // Assuming string representation for byte16
}

// Preconf bid structure
// PreConfBid represents the bid data.
type PreConfBid struct {
	UnsignedPreConfBid

	BidHash   []byte `json:"bidhash"`
	Signature []byte `json:"signature"`
}

// golang interface
type IPreconfBid interface {
	GetTxnHash() string
	GetBidAmt() *big.Int
	VerifySignature() (common.Address, error)
}

type IPreconfBidSearcher interface {
	IPreconfBid
	SubmitBid() error
}

// Returns a PreConfBid Object with an EIP712 signature of the payload
func ConstructSignedBid(bidamt *big.Int, txnhash string, key *ecdsa.PrivateKey) (IPreconfBid, error) {
	bid := &PreConfBid{
		UnsignedPreConfBid: UnsignedPreConfBid{
			Bid:     bidamt,
			TxnHash: txnhash,
		},
	}

	err := bid.constructHashAndSignature(key)
	if err != nil {
		return nil, err
	}

	return bid, nil
}

// Verifies the bid
func (p PreConfBid) GetTxnHash() string {
	return p.TxnHash
}

func (p PreConfBid) GetBidAmt() *big.Int {
	return p.Bid
}

// Verifies the bid
func (p PreConfBid) VerifySignature() (common.Address, error) {
	if p.BidHash == nil || p.Signature == nil {
		return common.Address{}, ErrMissingHashSignature
	}

	internalPayload := constructBidPayload(p.TxnHash, p.Bid)

	bidHash, _, err := encodeBidPayload(internalPayload)
	if err != nil {
		return common.Address{}, err
	}

	if !bytes.Equal(bidHash, p.BidHash) {
		return common.Address{}, ErrInvalidHash
	}

	pubkey, err := crypto.SigToPub(bidHash, p.Signature)
	if err != nil {
		return common.Address{}, err
	}

	if !crypto.VerifySignature(crypto.FromECDSAPub(pubkey), bidHash, p.Signature[:len(p.Signature)-1]) {
		return common.Address{}, ErrInvalidSignature
	}

	return crypto.PubkeyToAddress(*pubkey), err
}

// Adds bidHash and Signature to preconfbid
// Fails atomically
func (p *PreConfBid) constructHashAndSignature(privKey *ecdsa.PrivateKey) (err error) {
	if p.BidHash != nil || p.Signature != nil {
		return ErrAlreadySignedBid
	}

	internalPayload := constructBidPayload(p.TxnHash, p.Bid)

	bidHash, _, err := encodeBidPayload(internalPayload)
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(bidHash, privKey)
	if err != nil {
		return err
	}

	p.BidHash = bidHash
	p.Signature = sig

	return nil
}

// Constructs the EIP712 formatted bid
func constructBidPayload(txnHash string, bid *big.Int) apitypes.TypedData {
	signerData := apitypes.TypedData{
		Types: apitypes.Types{
			"PreConfBid": []apitypes.Type{
				{Name: "TxnHash", Type: "string"},
				{Name: "bid", Type: "uint64"},
			},
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
			},
		},
		PrimaryType: "PreConfBid",
		Domain: apitypes.TypedDataDomain{
			Name:    "PreConfBid",
			Version: "1",
		},
		Message: apitypes.TypedDataMessage{
			"TxnHash": txnHash,
			"bid":     bid,
		},
	}

	return signerData
}

func encodeBidPayload(preConfBid apitypes.TypedData) ([]byte, string, error) {
	return apitypes.TypedDataAndHash(preConfBid)
}
