package preconf_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"os"
	"testing"

	builder "github.com/primev/builder-boost/pkg/builder/preconf"
	"github.com/primev/builder-boost/pkg/preconf"
	"golang.org/x/exp/slog"
)

type testBidReader struct {
}

func (r *testBidReader) BidReader() <-chan preconf.PreConfBid {
	return nil
}

type testCommitmentWriter struct {
}

func (w *testCommitmentWriter) WriteCommitment(commitment preconf.PreconfCommitment) error {
	return nil
}

func newTestLogger() *slog.Logger {
	testLogger := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(testLogger)
}

func TestPreconfWorker(t *testing.T) {
	t.Parallel()

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	reader := &testBidReader{}
	writer := &testCommitmentWriter{}

	worker := builder.NewPreconfWorker(reader, writer, newTestLogger(), privateKey)

	t.Cleanup(func() {
		err := worker.Close()
		if err != nil {
			t.Fatal(err)
		}
	})

}
