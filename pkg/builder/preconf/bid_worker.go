package preconf

import (
	"crypto/ecdsa"
	"io"

	"github.com/primev/builder-boost/pkg/preconf"
	"golang.org/x/exp/slog"
)

type BidReader interface {
	BidReader() <-chan preconf.PreConfBid
}

type CommitmentWriter interface {
	WriteCommitment(commitment preconf.PreconfCommitment) error
}

type PreconfWorker struct {
	io.Closer

	reader     BidReader
	writer     CommitmentWriter
	logger     *slog.Logger
	builderKey *ecdsa.PrivateKey
	quit       chan struct{}
}

func NewPreconfWorker(
	reader BidReader,
	writer CommitmentWriter,
	logger *slog.Logger,
	builderKey *ecdsa.PrivateKey,
) *PreconfWorker {
	p := &PreconfWorker{
		reader:     reader,
		writer:     writer,
		logger:     logger,
		builderKey: builderKey,
		quit:       make(chan struct{}),
	}

	go p.run()
	return p
}

func (w *PreconfWorker) run() {
	for {
		select {
		case <-w.quit:
			w.logger.Info("preconf worker quit")
			return
		case bid, more := <-w.reader.BidReader():
			if !more {
				w.logger.Info("preconf bid reader closed, closing preconf worker")
				return
			}
			address, err := bid.VerifySearcherSignature()
			if err != nil {
				w.logger.Error("failed to verify searcher signature", "err", err)
				continue
			}

			w.logger.Info("preconf bid verified",
				"address", address.Hex(),
				"bid_tnx", bid.TxnHash,
				"bid_amt", bid.GetBidAmt(),
			)

			commitment, err := bid.ConstructCommitment(w.builderKey)
			if err != nil {
				w.logger.Error("failed to construct commitment", "err", err)
				continue
			}

			w.logger.Info("commitment constructed", "commitment", commitment)

			err = w.writer.WriteCommitment(commitment)
			if err != nil {
				w.logger.Error("failed to write commitment", "err", err)
			}
		}
	}
}

func (w *PreconfWorker) Close() error {
	close(w.quit)
	return nil
}
