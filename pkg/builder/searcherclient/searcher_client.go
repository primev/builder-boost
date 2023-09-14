package searcherclient

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"math/big"
	"sync"
	"time"

	"github.com/attestantio/go-builder-client/api/capella"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/exp/slog"
)

// SearcherInfo is a struct containing information about a searcher
type SearcherInfo struct {
	ID        string    `json:"id"`
	Validity  int64     `json:"validity"`
	Heartbeat time.Time `json:"heartbeat"`
}

// SearcherClient is an interface for managing searcher connections
type SearcherClient interface {
	// AddSearcher adds a searcher to the client
	AddSearcher(s Searcher)
	// RemoveSearcher removes a searcher from the client
	RemoveSearcher(id string)
	// Disconnected is called when a searcher disconnects
	Disconnected(id, reason string)
	// IsConnected returns true if the searcher is connected
	IsConnected(id string) bool
	// Send sends a payload to all relevant searchers
	SubmitBlock(ctx context.Context, payload *capella.SubmitBlockRequest) error
	// GetSearcherInfo returns a list of searcher info
	GetSeacherInfo() []SearcherInfo

	io.Closer
}

type Transaction struct {
	Count          int64 `json:"count"`
	MinPriorityFee int64 `json:"MinPriorityFee"`
	MaxPriorityFee int64 `json:"MaxPriorityFee"`
}

type Metadata struct {
	Builder            string      `json:"builder"`
	Number             int64       `json:"number"`
	BlockHash          string      `json:"blockHash"`
	Timestamp          string      `json:"timestamp"`
	BaseFee            uint32      `json:"baseFee"`
	Transactions       Transaction `json:"standard_transactions"`
	ClientTransactions []string    `json:"personal_transactions,omitempty"`
	SentTimestamp      time.Time   `json:"sent_timestamp"` // Timestamp of block sent to the searcher
	RecTimestamp       time.Time   `json:"rec_timestamp"`  // Timestamp of block received by the builder instance
}

type SuperPayload struct {
	InternalMetadata Metadata
	SearcherTxns     map[string][]string
}

type searcherClient struct {
	mu                   sync.Mutex
	searchers            map[string]Searcher
	disconnected         map[string]string
	inclusionProofActive bool
	logger               *slog.Logger
}

func NewSearcherClient(logger *slog.Logger, inclusionProof bool) SearcherClient {
	return &searcherClient{
		searchers:            make(map[string]Searcher),
		logger:               logger,
		inclusionProofActive: inclusionProof,
	}
}

func (s *searcherClient) AddSearcher(searcher Searcher) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.searchers[searcher.ID()] = searcher
	delete(s.disconnected, searcher.ID())
}

func (s *searcherClient) RemoveSearcher(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	searcher, found := s.searchers[id]
	if !found {
		return
	}
	err := searcher.Close()
	if err != nil {
		s.logger.Error("failed to close searcher", "err", err)
	}
	delete(s.searchers, searcher.ID())
}

func (s *searcherClient) Disconnected(id, reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.disconnected[id] = reason
	delete(s.searchers, id)
}

func (s *searcherClient) SubmitBlock(
	ctx context.Context,
	payload *capella.SubmitBlockRequest,
) error {
	if payload == nil || len(payload.ExecutionPayload.Transactions) == 0 {
		return nil
	}

	minTipTxn, maxTipTxn := big.NewInt(0), big.NewInt(0)
	ts := time.Unix(int64(payload.ExecutionPayload.Timestamp), 0).Format(time.RFC1123)
	baseFee := binary.LittleEndian.Uint32(payload.ExecutionPayload.BaseFeePerGas[:])

	blockMetadata := SuperPayload{
		InternalMetadata: Metadata{
			RecTimestamp: time.Now(),
			Builder:      payload.Message.BuilderPubkey.String(),
			Number:       int64(payload.ExecutionPayload.BlockNumber),
			BlockHash:    payload.Message.BlockHash.String(),
			Timestamp:    ts,
			BaseFee:      baseFee,
			Transactions: Transaction{
				Count: int64(len(payload.ExecutionPayload.Transactions)),
			},
		},
	}
	if s.inclusionProofActive {
		blockMetadata.SearcherTxns = make(map[string][]string)
	}

	for idx, btxn := range payload.ExecutionPayload.Transactions {
		var txn types.Transaction
		err := txn.UnmarshalBinary(btxn)
		if err != nil {
			s.logger.Error("failed to decode transaction", "err", err)
			continue
		}
		// Extract Min/Max
		if txn.GasTipCap().Cmp(minTipTxn) < 0 || idx == 0 {
			minTipTxn = txn.GasTipCap()
		}
		if txn.GasTipCap().Cmp(maxTipTxn) > 0 || idx == 0 {
			maxTipTxn = txn.GasTipCap()
		}

		from, err := types.Sender(types.LatestSignerForChainID(txn.ChainId()), &txn)
		if err != nil {
			s.logger.Error("failed to decode sender of transaction",
				"err", err,
				"txn", txn.Hash().String(),
			)
			continue
		}
		clientID := from.Hex()
		if s.inclusionProofActive {
			blockMetadata.SearcherTxns[clientID] = append(
				blockMetadata.SearcherTxns[clientID],
				txn.Hash().String(),
			)
		}
	}

	blockMetadata.InternalMetadata.Transactions.MinPriorityFee = minTipTxn.Int64()
	blockMetadata.InternalMetadata.Transactions.MaxPriorityFee = maxTipTxn.Int64()

	s.logger.Info("block metadata processed", "blockMetadata", blockMetadata)

	if s.inclusionProofActive {
		s.mu.Lock()
		defer s.mu.Unlock()

		for _, searcher := range s.searchers {
			if _, ok := blockMetadata.SearcherTxns[searcher.ID()]; ok {
				err := searcher.Send(ctx, blockMetadata)
				if err != nil {
					s.logger.Error("failed to send block metadata to searcher",
						"err", err,
						"searcher", searcher.ID(),
					)
				}
				// this could only happen if the callback for disconnection is
				// not called yet
			}
		}
	}

	return nil
}

func (s *searcherClient) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	for _, searcher := range s.searchers {
		err = errors.Join(searcher.Close(), err)
		delete(s.searchers, searcher.ID())
	}

	return err
}

func (s *searcherClient) IsConnected(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.searchers[id]
	return exists
}

func (s *searcherClient) GetSeacherInfo() []SearcherInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	searcherInfo := make([]SearcherInfo, 0, len(s.searchers))
	for _, searcher := range s.searchers {
		searcherInfo = append(searcherInfo, searcher.GetInfo())
	}
	return searcherInfo
}
