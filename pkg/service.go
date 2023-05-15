package boost

import (
	"context"
	"sync"
	"time"

	"github.com/attestantio/go-builder-client/api/capella"
	"github.com/lthibault/log"
)

const (
	Version = "0.0.1"
)

type BoostService interface {
	// Primev APIs
	// Example(context.Context, *types.ExRequest) error
	SubmitBlock(context.Context, *capella.SubmitBlockRequest) error
}

type DefaultService struct {
	Log    log.Logger
	Config Config
	Boost

	once  sync.Once
	ready chan struct{}
}

func (s *DefaultService) Run(ctx context.Context) (err error) {
	if s.Log == nil {
		s.Log = log.New().WithField("service", "boost")
	}

	timeBoostStart := time.Now()
	if s.Boost == nil {
		s.Boost, err = NewBoost(s.Config)
		if err != nil {
			return
		}
	}
	s.Log.With(log.F{
		"service":     "boost",
		"startTimeMs": time.Since(timeBoostStart).Milliseconds(),
	}).Info("initialized")

	s.setReady()

	return nil
}

func (s *DefaultService) Ready() <-chan struct{} {
	s.once.Do(func() {
		s.ready = make(chan struct{})
	})
	return s.ready
}

func (s *DefaultService) setReady() {
	select {
	case <-s.Ready():
	default:
		close(s.ready)
	}
}
