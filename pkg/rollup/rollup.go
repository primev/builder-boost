package rollup

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

//go:generate mockery --name Rollup

type Rollup interface {
	// Run starts rollup contract event listener
	Run(ctx context.Context) error

	// GetBuilderAddress returns current builder address
	GetBuilderAddress() common.Address

	// GetStake returns cached stake of searcher commited to specified builder
	GetStake(searcher common.Address, commitment common.Hash) *big.Int

	// GetAggregaredStake returns aggregated cached stake by searcher address
	GetAggregaredStake(searcher common.Address) *big.Int

	// GetStakeRemote fetches and returns balance of searcher commited to this builder from remote contract.
	// After fetching result is cached.
	GetStakeRemote(searcher common.Address, commitment common.Hash) (*big.Int, error)

	// GetMinimalStake returns cached minimal stake of specified builder
	GetMinimalStake(builder common.Address) *big.Int

	// IsSyncing returns true if service is in sync state
	IsSyncing() bool
}

func New(
	client *ethclient.Client,
	contractAddress common.Address,
	builderKey *ecdsa.PrivateKey,
	startBlock uint64,
	statePath string,
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
		startBlock:     startBlock,

		statePath:    statePath,
		state:        State{},
		stateMutex:   sync.Mutex{},
		stateUpdated: false,

		log: log.WithField("service", "rollup"),
	}

	// load rollup state
	err = r.loadState()
	if err != nil {
		return nil, err
	}

	return r, nil
}

type rollup struct {
	client   *ethclient.Client
	contract *contracts.BuilderStaking

	builderKey     *ecdsa.PrivateKey
	builderAddress common.Address
	startBlock     uint64

	statePath    string
	state        State
	stateMutex   sync.Mutex
	stateUpdated bool

	log log.Logger
}

// Run starts rollup contract event listener
func (r *rollup) Run(ctx context.Context) error {
	// process events from rollup contract
	for {
		blockProcessTimer := time.After(blockProcessPeriod)

		// process events from next batch of blocks
		err := r.processNextBlocks(ctx)
		if err != nil {
			return err
		}

		// save rollup state after processing events
		err = r.saveState()
		if err != nil {
			return err
		}

		// delay processing new batch of blocks to meet RPC rate limits
		<-blockProcessTimer
	}
}

// GetBuilderAddress returns current builder address
func (r *rollup) GetBuilderAddress() common.Address {
	return r.builderAddress
}

// GetStake returns cached stake of searcher commited to specified builder
func (r *rollup) GetStake(searcher common.Address, commitment common.Hash) *big.Int {
	return r.getStake(searcher, commitment)
}

// GetAggregaredStake returns aggregated cached stake by commitment hash
func (r *rollup) GetAggregaredStake(searcher common.Address) *big.Int {
	commitment := utils.GetCommitment(r.builderKey, searcher)
	return r.getAggregatedStake(commitment)
}

// GetStakeRemote fetches and returns balance of searcher commited to this builder from remote contract.
// After fetching result is cached.
func (r *rollup) GetStakeRemote(searcher common.Address, commitment common.Hash) (*big.Int, error) {
	stake, err := r.contract.Stakes(nil, searcher, commitment)
	if err != nil {
		return nil, err
	}

	r.setStake(searcher, commitment, stake)

	return stake, nil
}

// GetMinimalStake returns cached minimal stake of specified builder
func (r *rollup) GetMinimalStake(builder common.Address) *big.Int {
	return r.getMinimalStake(builder)
}

// IsSyncing returns true if service is still synced rollup
func (r *rollup) IsSyncing() bool {
	r.stateMutex.Lock()
	defer r.stateMutex.Unlock()

	return !r.stateUpdated || r.state.LatestKnownBlock > r.state.LatestProcessedBlock+blockSyncStateThreshold
}

// processNextEvents processes events from next batch of blocks and updates local state
func (r *rollup) processNextBlocks(ctx context.Context) error {
	// receive start and end block to process
	startBlock := r.state.LatestProcessedBlock + 1
	latestBlock, err := r.client.BlockNumber(ctx)
	if err != nil {
		return err
	}

	endBlock := latestBlock
	if startBlock+blockProcessBatchSize < latestBlock {
		endBlock = startBlock + blockProcessBatchSize
	}

	// no new blocks to process
	if startBlock > endBlock {
		return nil
	}

	if startBlock == endBlock {
		r.log.WithField("block", startBlock).Info("processing rollup block")
	} else {
		r.log.WithField("from", startBlock).WithField("to", endBlock).WithField("batch", endBlock-startBlock+1).
			Info("processing old rollup blocks in batch")
	}

	// process minimal stake updated events
	minimalStakeUpdatedIterator, err := r.contract.FilterMinimalStakeUpdated(&bind.FilterOpts{Start: startBlock, End: &endBlock, Context: ctx})
	if err != nil {
		return err
	}

	err = r.processMinimalStakeUpdatedEvents(minimalStakeUpdatedIterator)
	if err != nil {
		return err
	}

	// process stake updated events
	stakeUpdatedIterator, err := r.contract.FilterStakeUpdated(&bind.FilterOpts{Start: startBlock, End: &endBlock, Context: ctx})
	if err != nil {
		return err
	}

	err = r.processStakeUpdatedEvents(stakeUpdatedIterator)
	if err != nil {
		return err
	}

	r.stateMutex.Lock()
	defer r.stateMutex.Unlock()
	r.state.LatestProcessedBlock = endBlock
	r.state.LatestKnownBlock = latestBlock
	r.stateUpdated = true

	return nil
}

// processMinimalStakeUpdatedEvents processes minimal stake updated events and saves data to state
func (r *rollup) processMinimalStakeUpdatedEvents(it *contracts.BuilderStakingMinimalStakeUpdatedIterator) error {
	for it.Next() {
		if it.Error() != nil {
			return it.Error()
		}

		blockNumber := it.Event.Raw.BlockNumber
		builder := it.Event.Builder
		minimalStake := it.Event.MinimalStake

		r.setMinimalStake(builder, minimalStake)

		r.log.WithField("builder", builder).WithField("minimalStake", minimalStake).WithField("block", blockNumber).
			Info("processed minimal stake updated event")
	}

	return nil
}

// processStakeUpdatedEvents processes stake updated events and saves data to state
func (r *rollup) processStakeUpdatedEvents(it *contracts.BuilderStakingStakeUpdatedIterator) error {
	for it.Next() {
		if it.Error() != nil {
			return it.Error()
		}

		blockNumber := it.Event.Raw.BlockNumber
		commitment := common.Hash(it.Event.Commitment)
		searcher := it.Event.Searcher
		stake := it.Event.Stake

		r.setStake(searcher, commitment, stake)

		r.log.WithField("searcher", searcher).WithField("commitment", commitment).
			WithField("stake", stake).WithField("block", blockNumber).
			Info("processed stake updated event")
	}

	return nil
}

// loadState reads state from state file and stores it in local rollup instance
func (r *rollup) loadState() error {
	// create and save clean state if file does not exists
	if !r.fileExists(r.statePath) {
		r.state = State{
			LatestProcessedBlock: r.startBlock,
			LatestKnownBlock:     r.startBlock,
			Stakes:               make(map[common.Address]map[common.Hash]BigInt),
			AggregatedStakes:     make(map[common.Hash]BigInt),
			MinimalStakes:        make(map[common.Address]BigInt),
		}

		return r.saveState()
	}

	stateFile, err := os.Open(r.statePath)
	if err != nil {
		return err
	}
	defer stateFile.Close()

	stateBytes, err := ioutil.ReadAll(stateFile)
	if err != nil {
		return err
	}

	var state State
	err = json.Unmarshal(stateBytes, &state)
	if err != nil {
		return err
	}

	if state.Stakes == nil {
		state.Stakes = make(map[common.Address]map[common.Hash]BigInt)
	}

	if state.MinimalStakes == nil {
		state.MinimalStakes = make(map[common.Address]BigInt)
	}

	r.state = state

	return nil
}

// saveState saves local rollup state to state file
func (r *rollup) saveState() error {
	stateFile, err := os.OpenFile(r.statePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer stateFile.Close()

	encoder := json.NewEncoder(stateFile)
	return encoder.Encode(r.state)
}

// getStake returns stake value staked for particular builder by searcher
func (r *rollup) getStake(searcher common.Address, commitment common.Hash) *big.Int {
	r.stateMutex.Lock()
	defer r.stateMutex.Unlock()

	searcherStakes, ok := r.state.Stakes[searcher]
	if !ok {
		return big.NewInt(0)
	}

	stake, ok := searcherStakes[commitment]
	if !ok {
		return big.NewInt(0)
	}

	return &stake.Int
}

// getAggregatedStake returns aggregated stake value staked for commitment
func (r *rollup) getAggregatedStake(commitment common.Hash) *big.Int {
	r.stateMutex.Lock()
	defer r.stateMutex.Unlock()

	stake, ok := r.state.AggregatedStakes[commitment]
	if !ok {
		return big.NewInt(0)
	}

	return &stake.Int
}

// setStake updates stake value staked for particular builder by searcher
func (r *rollup) setStake(searcher common.Address, commitment common.Hash, stake *big.Int) {
	r.stateMutex.Lock()
	defer r.stateMutex.Unlock()

	if _, ok := r.state.Stakes[searcher]; !ok {
		r.state.Stakes[searcher] = make(map[common.Hash]BigInt)
	}

	oldStake := r.state.Stakes[searcher][commitment].Int
	oldAggregatedStake := r.state.AggregatedStakes[commitment].Int
	diff := stake.Sub(stake, &oldStake)

	r.state.Stakes[searcher][commitment] = BigInt{*stake}
	r.state.AggregatedStakes[commitment] = BigInt{*big.NewInt(0).Add(&oldAggregatedStake, diff)}
}

// getMinimalStake returns stake value staked for particular builder by searcher
func (r *rollup) getMinimalStake(builder common.Address) *big.Int {
	r.stateMutex.Lock()
	defer r.stateMutex.Unlock()

	stake, ok := r.state.MinimalStakes[builder]
	if !ok {
		return big.NewInt(0)
	}

	return &stake.Int
}

// setMinimalStake updates minimal stake value for particular builder
func (r *rollup) setMinimalStake(builder common.Address, stake *big.Int) {
	r.stateMutex.Lock()
	defer r.stateMutex.Unlock()

	r.state.MinimalStakes[builder] = BigInt{*stake}
}

// fileExists returns true if file under specified file path exists
func (r *rollup) fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}
