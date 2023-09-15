package searcherclient

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

var wsTimeout = 10 * time.Second

// Searcher is the interface that wraps the lifecycle of a searcher connection.
type Searcher interface {
	// ID returns the ID of the searcher.
	ID() string
	// Send sends a payload to the searcher.
	Send(ctx context.Context, payload SuperPayload) error
	// GetInfo returns information about the searcher.
	GetInfo() SearcherInfo

	io.Closer
}

// RollUp is the interface that wraps the block number of a rollup contract.
type RollUp interface {
	GetBlockNumber() (*big.Int, error)
}

// DisconnectHandler is a function that is called when a searcher disconnects.
type DisconnectHandler func(reason string)

var (
	// ErrQuit is returned when user issues a close.
	ErrQuit = errors.New("server quit")
)

// NewSearcher creates a new searcher connection.
func NewSearcher(
	id string,
	conn *websocket.Conn,
	rollup RollUp,
	logger *slog.Logger,
	subscriptionEnd int64,
	disconnectCB DisconnectHandler,
) Searcher {
	s := &searcherConn{
		wsConn:          conn,
		searcherID:      id,
		quit:            make(chan struct{}),
		outbox:          make(chan wsMsg),
		logger:          logger,
		rollup:          rollup,
		subscriptionEnd: subscriptionEnd,
		disconnectCB:    disconnectCB,
	}

	s.heartbeat.Store(time.Now().Unix())

	go s.startWorker()
	return s
}

type wsMsg struct {
	msgType int
	msg     []byte
}

type searcherConn struct {
	wsConn          *websocket.Conn
	searcherID      string
	quit            chan struct{}
	subscriptionEnd int64
	logger          *slog.Logger
	heartbeat       atomic.Int64
	outbox          chan wsMsg
	rollup          RollUp
	disconnectCB    DisconnectHandler
}

func (s *searcherConn) isSubscriptionEnded() bool {
	blkNo, err := s.rollup.GetBlockNumber()
	if err != nil {
		// if we can't get the block number, assume the server is not in a good state
		// so we should disconnect
		s.logger.Error("failed to get block number", "err", err)
		return true
	}
	return blkNo.Int64() > s.subscriptionEnd
}

func (s *searcherConn) startWorker() {
	for {
		select {
		case <-s.quit:
			s.disconnectCB("server closing connection")
			_ = s.wsConn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(
					websocket.CloseNormalClosure,
					"closing connection",
				),
			)
			s.logger.Info("closing connection")
			return
		case payload := <-s.outbox:
			if s.isSubscriptionEnded() {
				s.disconnectCB("subscription ended")
				_ = s.wsConn.WriteMessage(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(
						websocket.ClosePolicyViolation,
						"subscription ended",
					),
				)
				s.logger.Info("subscription ended")
				return
			}
			err := s.wsConn.WriteMessage(payload.msgType, payload.msg)
			if err != nil {
				s.logger.Error("failed to send payload", "err", err)
				s.disconnectCB("failed to send payload")
				return
			}
			s.logger.Info("sent payload")
			s.heartbeat.Store(time.Now().Unix())
		}
	}
}

func (s *searcherConn) ID() string {
	return s.searcherID
}

func (s *searcherConn) GetInfo() SearcherInfo {
	return SearcherInfo{
		ID:        s.searcherID,
		Heartbeat: time.Unix(s.heartbeat.Load(), 0),
		Validity:  s.subscriptionEnd,
	}
}

func (s *searcherConn) Close() error {
	close(s.quit)
	return nil
}

func (s *searcherConn) Send(ctx context.Context, payload SuperPayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		s.logger.Error("failed to marshal payload", "err", err)
		return err
	}

	select {
	case <-s.quit:
		return ErrQuit
	default:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.outbox <- wsMsg{
		msgType: websocket.TextMessage,
		msg:     jsonPayload,
	}:
	}

	return nil
}
