package boost

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lthibault/log"
)

type SearcherConnection struct {
	Conn        *websocket.Conn
	WorkChannel chan SuperPayload
}

type Worker struct {
	lock               sync.RWMutex
	connectedSearchers map[string]SearcherConnection

	// Note: Heartbeat is meant to be accessed via atomic operations .Store and .Load to ensure non-blocking performance
	heartbeat *atomic.Int64

	log       log.Logger
	workQueue chan SuperPayload

	once  sync.Once
	ready chan struct{}
}

func (w *Worker) GetHeartbeat() int64 {
	return w.heartbeat.Load()
}

func NewWorker(workQueue chan SuperPayload, logger log.Logger) *Worker {
	return &Worker{
		connectedSearchers: make(map[string]SearcherConnection),
		workQueue:          workQueue,
		log:                logger,
		heartbeat:          &atomic.Int64{},
	}
}

// workerRecovery recovers from a panic in Worker.Run and restarts the Worker
func (w *Worker) workerRecovery(ctx context.Context) {
	if r := recover(); r != nil {
		w.log.Error("recovered from panic in Worker.Run", "error", r)

		// Restart the Worker with a new context
		newCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		w.log.Info("restarting Worker.Run")
		if err := w.Run(newCtx); err != nil {
			w.log.Error("failed to restart Worker.Run", "error", err)
		}
	}
}

// TODO(@ckartik): Add a channel to request health status of worker
func (w *Worker) Run(ctx context.Context) (err error) {
	w.log.Info("starting Worker.Run")
	go func() {
		defer w.workerRecovery(ctx)
		for {
			w.heartbeat.Store(time.Now().Unix())

			select {
			case <-ctx.Done():
				return
			case blockMetadata := <-w.workQueue:
				w.log.Info("received block metadata", "block", blockMetadata)
				w.lock.RLock()
				for _, searcher := range w.connectedSearchers {
					// NOTE: Risk of Worker Blocking here, if searcher is not reading from channel
					searcher.WorkChannel <- blockMetadata
				}
				w.lock.RUnlock()
			}
		}
	}()

	w.setReady()

	return nil
}

func (w *Worker) Ready() <-chan struct{} {
	w.once.Do(func() {
		w.ready = make(chan struct{})
	})
	return w.ready
}

func (w *Worker) setReady() {
	select {
	case <-w.Ready():
	default:
		close(w.ready)
	}
}
