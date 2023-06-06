package rollup

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lthibault/log"
	"github.com/primev/builder-boost/pkg/contracts"
	"github.com/primev/builder-boost/pkg/utils"
)

const (
	blockSyncStateThreshold = 4
	blockProcessBatchSize   = 2048
	blockProcessPeriod      = time.Second * 3
)

var (
	ErrNoMinimalStakeSet = errors.New("R0001: no minimal stake set, to fix set minimal stake for builder on payment contract")
)

//go:generate mockery --name Rollup --filename mockrollup.go

type Rollup interface {
	// Run starts rollup contract event listener
	Run(ctx context.Context) error

	// GetBuilderAddress returns current builder address
	GetBuilderAddress() common.Address

	// GetSubscriptionEnd returns subscription end for specified commitment
	GetSubscriptionEnd(commitment common.Hash) (*big.Int, error)

	// GetMinimalStake returns minimal stake of specified builder
	GetMinimalStake(builder common.Address) (*big.Int, error)

	// GetCommitment calculates commitment hash for this builder by searcher address
	GetCommitment(searcher common.Address) common.Hash

	// GetBlockNumber returns latest blocks number from rollup
	GetBlockNumber() (*big.Int, error)
}

func New(
	client *ethclient.Client,
	contractAddress common.Address,
	builderKey *ecdsa.PrivateKey,
	log log.Logger,
) (Rollup, error) {
	contract, err := contracts.NewBuilderStaking(contractAddress, client)
	if err != nil {
		return nil, err
	}

	r := &rollup{
		client:   client,
		contract: contract,

		builderKey:     builderKey,
		builderAddress: crypto.PubkeyToAddress(builderKey.PublicKey),

		log: log.WithField("service", "rollup"),
	}

	return r, nil
}

type rollup struct {
	client   *ethclient.Client
	contract *contracts.BuilderStaking

	builderKey     *ecdsa.PrivateKey
	builderAddress common.Address

	log log.Logger
}

// Run starts rollup contract event listener
func (r *rollup) Run(ctx context.Context) error {
	// process events from rollup contract
	for {
		blockProcessTimer := time.After(blockProcessPeriod)

		// check if minimal stake is set for current builder
		_, err := r.GetMinimalStake(r.builderAddress)
		if err != nil {
			r.log.WithField("builder", r.builderAddress).WithField("err", err.Error()).Error("minimal stake is not set")
		}

		// delay processing new batch of blocks to meet RPC rate limits
		<-blockProcessTimer
	}
}

// GetBuilderAddress returns current builder address
func (r *rollup) GetBuilderAddress() common.Address {
	return r.builderAddress
}

// GetSubscriptionEnd returns subscription end for specified commitment
func (r *rollup) GetSubscriptionEnd(commitment common.Hash) (*big.Int, error) {
	stake, err := r.contract.Stakes(nil, commitment)
	if err != nil {
		return nil, err
	}

	return stake.Stake, nil
}

// GetMinimalStake returns minimal stake for specified builder
func (r *rollup) GetMinimalStake(builder common.Address) (*big.Int, error) {
	info, err := r.contract.Builders(nil, builder)
	if err != nil {
		return nil, err
	}

	return info.MinimalStake, nil
}

// GetCommitment calculates commitment hash for this builder by searcher address
func (r *rollup) GetCommitment(searcher common.Address) common.Hash {
	return utils.GetCommitment(r.builderKey, searcher)
}

// GetBlockNumber returns latest blocks number from rollup
func (r *rollup) GetBlockNumber() (*big.Int, error) {
	blockNumber, err := r.client.BlockNumber(context.TODO())
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetUint64(blockNumber), nil
}
