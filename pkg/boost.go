package boost

import (
	"context"
	"encoding/binary"
	"errors"
	"time"

	"github.com/attestantio/go-builder-client/api/capella"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lthibault/log"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Boost interface {
	SubmitBlock(context.Context, *capella.SubmitBlockRequest) error
	GetWorkChannel() chan Metadata
}

type DefaultBoost struct {
	config      Config
	pushChannel chan Metadata
}

type Transaction struct {
	Count          int64 `json:"count"`
	MinPriorityFee int64 `json:"MinPriorityFee"`
	MaxPriorityFee int64 `json:"MaxPriorityFee"`
}

type Metadata struct {
	Builder      string      `json:"builder"`
	Number       int64       `json:"number"`
	BlockHash    string      `json:"blockHash"`
	Timestamp    string      `json:"timestamp"`
	BaseFee      uint32      `json:"baseFee"`
	Transactions Transaction `json:"transactions"`
}

var (
	ErrBlockUnprocessable = errors.New("V001: block unprocessable")
)

// NewGateway new auction gateway service
func NewBoost(config Config) (*DefaultBoost, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	as := &DefaultBoost{
		config:      config,
		pushChannel: make(chan Metadata, 100),
	}
	return as, nil
}

func (rs *DefaultBoost) GetWorkChannel() chan Metadata {
	return rs.pushChannel
}

func (rs *DefaultBoost) Log() log.Logger {
	return rs.config.Log
}

func (as *DefaultBoost) SubmitBlock(ctx context.Context, msg *capella.SubmitBlockRequest) (err error) {
	span, _ := tracer.StartSpanFromContext(ctx, "submit-block")
	defer span.Finish()
	defer func() {
		if r := recover(); r != nil {
			as.config.Log.Error(ErrBlockUnprocessable.Error())
			err = ErrBlockUnprocessable
		}
	}()

	var _txn types.Transaction
	var blockMetadata Metadata

	blockMetadata.BlockHash = msg.Message.BlockHash.String()
	blockMetadata.Number = int64(msg.ExecutionPayload.BlockNumber)
	blockMetadata.Builder = msg.Message.BuilderPubkey.String()
	blockMetadata.Transactions.Count = int64(len(msg.ExecutionPayload.Transactions))
	blockMetadata.Timestamp = time.Unix(int64(msg.ExecutionPayload.Timestamp), 0).Format(time.RFC1123)
	blockMetadata.BaseFee = binary.LittleEndian.Uint32(msg.ExecutionPayload.BaseFeePerGas[:])

	// Conditionally set txn details based on txn count
	if len(msg.ExecutionPayload.Transactions) > 0 {
		_txn.UnmarshalBinary(msg.ExecutionPayload.Transactions[0])
		minTipTxn := _txn
		maxTipTxn := _txn
		for _, btxn := range msg.ExecutionPayload.Transactions {
			var txn types.Transaction
			txn.UnmarshalBinary(btxn)
			// Extract Min/Max
			if txn.GasTipCapCmp(&minTipTxn) < 0 {
				minTipTxn = txn
			}
			if txn.GasTipCapCmp(&maxTipTxn) > 0 {
				maxTipTxn = txn
			}
		}

		blockMetadata.Transactions.MinPriorityFee = minTipTxn.GasTipCap().Int64()
		blockMetadata.Transactions.MaxPriorityFee = maxTipTxn.GasTipCap().Int64()
	}

	as.pushChannel <- blockMetadata

	as.config.Log.
		WithField("block_hash", blockMetadata.BlockHash).
		WithField("base_fee", blockMetadata.BaseFee).
		WithField("min_priority_fee", blockMetadata.Transactions.MinPriorityFee).
		WithField("max_priority_fee", blockMetadata.Transactions.MaxPriorityFee).
		WithField("txn_count", blockMetadata.Transactions.Count).
		WithField("builder", blockMetadata.Builder).
		Info("Block metadata processed")

	return nil
}
