package preconf

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/primev/builder-boost/pkg/contracts"
	"github.com/primev/builder-boost/pkg/p2p/commons"
	"github.com/primev/builder-boost/pkg/p2p/node"
)

var (
	ErrAlreadySignedBid        = errors.New("already contains hash or signature")
	ErrMissingHashSignature    = errors.New("missing hash or signature")
	ErrInvalidSignature        = errors.New("signature is not valid")
	ErrInvalidHash             = errors.New("bidhash doesn't match bid payload")
	ErrAlreadySignedCommitment = errors.New("commitment is already hashed or signed")
)

type UnsignedPreConfBid struct {
	TxnHash     string   `json:"txnHash"`
	Bid         *big.Int `json:"bid"`
	Blocknumber *big.Int `json:"blocknumber"`
	// UUID    string `json:"uuid"` // Assuming string representation for byte16
}

// Most of the bid details can go to the data-availabilty layer if needed, and only have the hash+sig live on chian
// Preconf bid structure
// PreConfBid represents the bid data.
type PreConfBid struct { // Adds blocknumber for pre-conf bid - Will need to manage how to reciever acts on a bid / TTL is the blocknumber
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
	BidOriginator() (common.Address, *ecdsa.PublicKey, error)
}

type IPreconfCommitment interface {
	IPreconfBid // TODO Implement for underlying data-structure
	VerifyBuilderSignature() (common.Address, error)
	CommitmentOriginator() (common.Address, *ecdsa.PublicKey, error)
}

type IPreconfCommitmentBuilder interface {
	IPreconfCommitment
	PublishCommitment() error
	StoreCommitmentToDA(*ecdsa.PrivateKey, *ethclient.Client) (*types.Transaction, error) // // TODO(@ckartik): Turn into Singleton client for production
}

type IPreconfBidSearcher interface {
	IPreconfBid
	SubmitBid(node.ISearcherNode) error
}

type IPreconfBidBuilder interface {
	IPreconfBid
	ConstructCommitment(*ecdsa.PrivateKey) (PreconfCommitment, error) // Verfiy Signature and than constrcut the commitment
}

func (p PreconfCommitment) StoreCommitmentToDA(privateKey *ecdsa.PrivateKey, contractAddress string, client *ethclient.Client) (*types.Transaction, error) {
	preconf, err := contracts.NewPreConfCommitmentStore(common.HexToAddress(contractAddress), client)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to retrieve nonce: %v", err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	txn, err := preconf.StoreCommitment(auth, p.PreConfBid.TxnHash, p.Bid.Uint64(), p.Blocknumber.Uint64(), hex.EncodeToString(p.BidHash), hex.EncodeToString(p.PreConfBid.Signature), hex.EncodeToString(p.CommitmentSignature))
	if err != nil {
		return nil, err
	}
	return txn, nil
}

func (p PreconfCommitment) VerifyBuilderSignature() (common.Address, error) {
	if p.DataHash == nil || p.CommitmentSignature == nil {
		return common.Address{}, ErrMissingHashSignature
	}

	internalPayload := constructCommitmentPayload(p.TxnHash, p.Bid, p.Blocknumber, p.BidHash, p.Signature)

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

func (p PreConfBid) BidOriginator() (common.Address, *ecdsa.PublicKey, error) {
	_, err := p.VerifySearcherSignature()
	if err != nil {
		return common.Address{}, nil, err
	}

	pubkey, err := crypto.SigToPub(p.BidHash, p.Signature)
	if err != nil {
		return common.Address{}, nil, err
	}

	return crypto.PubkeyToAddress(*pubkey), pubkey, nil
}

func (p PreconfCommitment) CommitmentOriginator() (common.Address, *ecdsa.PublicKey, error) {
	_, err := p.VerifyBuilderSignature()
	if err != nil {
		return common.Address{}, nil, err
	}

	pubkey, err := crypto.SigToPub(p.DataHash, p.CommitmentSignature)
	if err != nil {
		return common.Address{}, nil, err
	}

	return crypto.PubkeyToAddress(*pubkey), pubkey, nil
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

	eip712Payload := constructCommitmentPayload(p.TxnHash, p.Bid, p.Blocknumber, p.BidHash, p.Signature)

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

func (p PreConfBid) SubmitBid(p2pEngine node.ISearcherNode) error {
	// searlize p into json
	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}

	fmt.Println(payload)
	p2pEngine.BidSend(commons.Publish, payload)
	// send payload to server at localhost:8080/preconf
	//http client
	request := http.Client{}
	request.Post("http://localhost:8080/preconf", "application/json", bytes.NewBuffer(payload))

	return nil
}

// Returns a PreConfBid Object with an EIP712 signature of the payload
func ConstructSignedBid(bidamt *big.Int, txnhash string, blocknumber *big.Int, key *ecdsa.PrivateKey) (IPreconfBidSearcher, error) {
	bid := &PreConfBid{
		UnsignedPreConfBid: UnsignedPreConfBid{
			Bid:         bidamt,
			TxnHash:     txnhash,
			Blocknumber: blocknumber,
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

	return eipVerify(constructBidPayload(p.TxnHash, p.Bid, p.Blocknumber), p.BidHash, p.Signature)
}

// Adds bidHash and Signature to preconfbid
// Fails atomically
func (p *PreConfBid) constructHashAndSignature(privKey *ecdsa.PrivateKey) (err error) {
	if p.BidHash != nil || p.Signature != nil {
		return ErrAlreadySignedBid
	}

	internalPayload := constructBidPayload(p.TxnHash, p.Bid, p.Blocknumber)

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
func constructCommitmentPayload(txnHash string, bid *big.Int, blockNumber *big.Int, bidHash []byte, signature []byte) apitypes.TypedData {
	signerData := apitypes.TypedData{
		Types: apitypes.Types{
			"PreConfCommitment": []apitypes.Type{
				{Name: "txnHash", Type: "string"},
				{Name: "bid", Type: "uint64"},
				{Name: "blockNumber", Type: "uint64"},
				{Name: "bidHash", Type: "string"},   // Hex Encoded Hash
				{Name: "signature", Type: "string"}, // Hex Encoded Signature
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
			"txnHash":     txnHash,
			"bid":         bid,
			"blockNumber": blockNumber,
			"bidHash":     hex.EncodeToString(bidHash),
			"signature":   hex.EncodeToString(signature),
		},
	}

	return signerData
}

// Constructs the EIP712 formatted bid
func constructBidPayload(txnHash string, bid *big.Int, blockNumber *big.Int) apitypes.TypedData {
	signerData := apitypes.TypedData{
		Types: apitypes.Types{
			"PreConfBid": []apitypes.Type{
				{Name: "txnHash", Type: "string"},
				{Name: "bid", Type: "uint64"},
				{Name: "blockNumber", Type: "uint64"},
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
			"txnHash":     txnHash,
			"bid":         bid,
			"blockNumber": blockNumber,
		},
	}

	return signerData
}
