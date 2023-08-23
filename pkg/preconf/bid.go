package preconf

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

var (
	ErrAlreadySignedBid        = errors.New("already contains hash or signature")
	ErrMissingHashSignature    = errors.New("missing hash or signature")
	ErrInvalidSignature        = errors.New("signature is not valid")
	ErrInvalidHash             = errors.New("bidhash doesn't match bid payload")
	ErrAlreadySignedCommitment = errors.New("commitment is already hashed or signed")
)

type UnsignedPreConfBid struct {
	TxnHash string   `json:"txnHash"`
	Bid     *big.Int `json:"bid"`
	// UUID    string `json:"uuid"` // Assuming string representation for byte16
}

// Most of the bid details can go to the data-availabilty layer if needed, and only have the hash+sig live on chian
// Preconf bid structure
// PreConfBid represents the bid data.
type PreConfBid struct {
	UnsignedPreConfBid

	BidHash   []byte `json:"bidhash"` // TODO(@ckaritk): name better
	Signature []byte `json:"signature"`
}

type PreconfCommitment struct {
	PreConfBid

	DataHash            []byte `json:"data_hash"` // TODO(@ckaritk): name better
	CommitmentSignature []byte `json:"commitment_signature"`
}

// golang interface
type IPreconfBid interface {
	GetTxnHash() string
	GetBidAmt() *big.Int
	VerifySearcherSignature() (common.Address, error)
}

type IPreconfCommitment interface {
	IPreconfBid
	VerifyBuilderSignature() (common.Address, error)
}

type IPreconfCommitmentBuilder interface {
	IPreconfCommitment
	PublishCommitment() error
}

type IPreconfBidSearcher interface {
	IPreconfBid
	SubmitBid() error
}

type IPreconfBidBuilder interface {
	IPreconfBid
	ConstructCommitment(*ecdsa.PrivateKey) (PreconfCommitment, error) // Verfiy Signature and than constrcut the commitment
}

// Verifies the bid
func (p PreconfCommitment) VerifyBuilderSignature() (common.Address, error) {
	if p.DataHash == nil || p.CommitmentSignature == nil {
		return common.Address{}, ErrMissingHashSignature
	}

	internalPayload := constructCommitmentPayload(p.TxnHash, p.Bid, p.BidHash, p.Signature)

	return eipVerify(internalPayload, p.DataHash, p.Signature)
}

func eipVerify(internalPayload apitypes.TypedData, expectedhash []byte, signature []byte) (common.Address, error) {
	payloadHash, _, err := apitypes.TypedDataAndHash(internalPayload)
	if err != nil {
		return common.Address{}, err
	}

	if !bytes.Equal(payloadHash, expectedhash) {
		return common.Address{}, ErrInvalidHash
	}

	pubkey, err := crypto.SigToPub(payloadHash, signature)
	if err != nil {
		return common.Address{}, err
	}

	if !crypto.VerifySignature(crypto.FromECDSAPub(pubkey), payloadHash, signature[:len(signature)-1]) {
		return common.Address{}, ErrInvalidSignature
	}

	return crypto.PubkeyToAddress(*pubkey), err
}

func (p PreConfBid) ConstructCommitment(privKey *ecdsa.PrivateKey) (PreconfCommitment, error) {
	_, err := p.VerifySearcherSignature()
	if err != nil {
		return PreconfCommitment{}, err
	}
	commitment := PreconfCommitment{
		PreConfBid: p,
	}

	err = commitment.constructHashAndSignature(privKey)
	if err != nil {
		return PreconfCommitment{}, err
	}

	return commitment, nil
}

// Adds bidHash and Signature to preconfbid
// Fails atomically
func (p *PreconfCommitment) constructHashAndSignature(privKey *ecdsa.PrivateKey) (err error) {
	if p.DataHash != nil || p.CommitmentSignature != nil {
		return ErrAlreadySignedCommitment
	}

	eip712Payload := constructCommitmentPayload(p.TxnHash, p.Bid, p.BidHash, p.Signature)

	dataHash, _, err := apitypes.TypedDataAndHash(eip712Payload)
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(dataHash, privKey)
	if err != nil {
		return err
	}

	p.DataHash = dataHash
	p.CommitmentSignature = sig

	return nil
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
func (p PreConfBid) VerifySearcherSignature() (common.Address, error) {
	if p.BidHash == nil || p.Signature == nil {
		return common.Address{}, ErrMissingHashSignature
	}

	return eipVerify(constructBidPayload(p.TxnHash, p.Bid), p.BidHash, p.Signature)
}

// Adds bidHash and Signature to preconfbid
// Fails atomically
func (p *PreConfBid) constructHashAndSignature(privKey *ecdsa.PrivateKey) (err error) {
	if p.BidHash != nil || p.Signature != nil {
		return ErrAlreadySignedBid
	}

	internalPayload := constructBidPayload(p.TxnHash, p.Bid)

	bidHash, _, err := apitypes.TypedDataAndHash(internalPayload)
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
func constructCommitmentPayload(txnHash string, bid *big.Int, bidHash []byte, signature []byte) apitypes.TypedData {
	signerData := apitypes.TypedData{
		Types: apitypes.Types{
			"PreConfBid": []apitypes.Type{
				{Name: "TxnHash", Type: "string"},
				{Name: "bid", Type: "uint64"},
				{Name: "BidHash", Type: "string"},   // Hex Encoded Hash
				{Name: "Signature", Type: "string"}, // Hex Encoded Signature
			},
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
			},
		},
		PrimaryType: "PreConfCommitment",
		Domain: apitypes.TypedDataDomain{
			Name:    "PreConfCommitment",
			Version: "1",
		},
		Message: apitypes.TypedDataMessage{
			"TxnHash":   txnHash,
			"bid":       bid,
			"bidHash":   hex.EncodeToString(bidHash),
			"signature": hex.EncodeToString(signature),
		},
	}

	return signerData
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
