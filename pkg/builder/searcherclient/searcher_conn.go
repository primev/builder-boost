package searcherclient

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
)

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
	// ErrSearcherDisconnected is returned when a searcher is disconnected.
	ErrSearcherDisconnected = errors.New("searcher disconnected")
	// ErrSearcherSubscriptionEnded is returned when a searcher's subscription has ended.
	ErrSearcherSubscriptionEnded = errors.New("searcher subscription ended")
	// ErrQuit is returned when user issues a close.
	ErrQuit = errors.New("server quit")
	// ErrIO is returned when there is an IO error.
	ErrIO = errors.New("IO error")
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
		inbox:           make(chan wsMsg),
		outbox:          make(chan wsMsg),
		logger:          logger,
		rollup:          rollup,
		subscriptionEnd: subscriptionEnd,
		heartbeat:       time.Now().Unix(),
		disconnectCB:    disconnectCB,
	}

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
	mu              sync.Mutex
	disconnected    bool
	subscriptionEnd int64
	logger          *slog.Logger
	heartbeat       int64
	inbox           chan wsMsg
	outbox          chan wsMsg
	rollup          RollUp
	disconnectCB    DisconnectHandler
}

func (s *searcherConn) isDisconnected() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.disconnected
}

func (s *searcherConn) setDisconnected() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.disconnected = true
}

func (s *searcherConn) isSubscriptionEnded() bool {
	blkNo, _ := s.rollup.GetBlockNumber()
	return blkNo.Int64() > s.subscriptionEnd
}

func (s *searcherConn) setHeartbeat() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.heartbeat = time.Now().Unix()
}

func (s *searcherConn) startWorker() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg, egCtx := errgroup.WithContext(ctx)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Start the reader
	eg.Go(func() error {
		for {
			select {
			case <-s.quit:
				return ErrQuit
			case <-egCtx.Done():
				return egCtx.Err()
			default:
			}
			msgType, msg, err := s.wsConn.ReadMessage()
			if err != nil {
				return errors.Join(ErrIO, err)
			}
			select {
			case s.inbox <- wsMsg{msgType: msgType, msg: msg}:
			case <-s.quit:
				return ErrQuit
			}
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-s.quit:
				return ErrQuit
			case msg := <-s.inbox:
				switch msg.msgType {
				case websocket.CloseMessage:
					return ErrSearcherDisconnected
				case websocket.PingMessage:
					err := s.wsConn.WriteMessage(websocket.PongMessage, []byte("pong"))
					if err != nil {
						s.logger.Error("failed to send pong", "err", err)
						return errors.Join(ErrIO, err)
					}
				case websocket.PongMessage:
				default:
					s.logger.Error("unknown message type", "msgType", msg.msgType)
					// non-fatal error
				}
				s.setHeartbeat()
			case payload := <-s.outbox:
				err := s.wsConn.WriteMessage(payload.msgType, payload.msg)
				if err != nil {
					s.logger.Error("failed to send payload", "err", err)
					return errors.Join(ErrIO, err)
				}
				s.setHeartbeat()
			case <-ticker.C:
				if s.isSubscriptionEnded() {
					err := s.wsConn.WriteMessage(
						websocket.CloseMessage,
						websocket.FormatCloseMessage(
							websocket.ClosePolicyViolation,
							"subscription ended",
						),
					)
					return errors.Join(ErrSearcherSubscriptionEnded, err)
				}

				if time.Now().Unix()-s.heartbeat > 10 {
					err := s.wsConn.WriteMessage(websocket.PingMessage, []byte("ping"))
					if err != nil {
						s.logger.Error("failed to send ping", "err", err)
						return errors.Join(ErrIO, err)
					}
				}
			}
		}
	})

	switch err := eg.Wait(); {
	case err == nil:
		// this should never happen
		s.logger.Error("worker exited with no error")
	case errors.Is(err, ErrQuit):
		s.logger.Info("worker stopped", "err", err)
	case errors.Is(err, ErrSearcherDisconnected):
		s.logger.Error("searcher disconnected", "err", err)
		s.disconnectCB(err.Error())
	case errors.Is(err, ErrSearcherSubscriptionEnded):
		s.logger.Error("searcher subscription ended", "err", err)
		s.disconnectCB(err.Error())
	default:
		s.logger.Error("could not reach searcher", "err", err)
		s.disconnectCB(err.Error())
	}
}

func (s *searcherConn) ID() string {
	return s.searcherID
}

func (s *searcherConn) GetInfo() SearcherInfo {
	return SearcherInfo{
		ID:        s.searcherID,
		Heartbeat: time.Unix(s.heartbeat, 0),
	}
}

func (s *searcherConn) Close() error {
	close(s.quit)
	return s.wsConn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing connection"),
	)
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
