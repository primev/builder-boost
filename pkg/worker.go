package boost

import (
	"context"
	"sync"

	"github.com/lthibault/log"
)

type Worker struct {
	// TODO(@ckartik): Wrap with RWLock
	lock               sync.RWMutex
	connectedSearchers map[string]chan Metadata
	log                log.Logger
	workQueue          chan Metadata

	once  sync.Once
	ready chan struct{}
}

func NewWorker(workQueue chan Metadata, logger log.Logger) *Worker {
	return &Worker{
		connectedSearchers: make(map[string]chan Metadata),
		workQueue:          workQueue,
		log:                logger,
	}
}

func (w *Worker) Run(ctx context.Context) (err error) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case blockMetadata := <-w.workQueue:
				w.log.Info("received block metadata", "block", blockMetadata)
				w.lock.RLock()
				for _, searcher := range w.connectedSearchers {
					// NOTE: Risk of Worker Blocking here, if searcher is not reading from channel
					searcher <- blockMetadata
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
