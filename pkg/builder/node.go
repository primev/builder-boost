package builder

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/builder-boost/pkg/apiserver"
	"github.com/primev/builder-boost/pkg/builder/api"
	"github.com/primev/builder-boost/pkg/builder/searcherclient"
	"github.com/primev/builder-boost/pkg/rollup"
	"golang.org/x/exp/slog"
)

// Options are the options for the builder node
type Options struct {
	Version              string
	Logger               *slog.Logger
	ServerAddr           string
	TracingEndpoint      string
	RollupKey            string
	RollupAddr           string
	RollupContract       string
	BuilderToken         string
	MetricsEnabled       bool
	DAContract           string
	DARPCAddr            string
	InclusionListEnabled bool
}

type Node struct {
	server *http.Server
	io.Closer
}

func NewNode(
	ctx context.Context,
	opts Options,
) (*Node, error) {
	rollupKeyBytes := common.FromHex(opts.RollupKey)
	rollupKey := crypto.ToECDSAUnsafe(rollupKeyBytes)

	rollupClient, err := ethclient.Dial(opts.RollupAddr)
	if err != nil {
		return nil, err
	}

	rollupContract := common.HexToAddress(opts.RollupContract)

	ru, err := rollup.New(rollupClient, rollupContract, rollupKey, nil)
	if err != nil {
		return nil, err
	}

	sclient := searcherclient.NewSearcherClient(opts.Logger, opts.InclusionListEnabled)

	srv := apiserver.New(opts.Version, opts.Logger)
	api.RegisterAPI(opts.BuilderToken, srv, ru, opts.Logger, sclient)

	server := &http.Server{
		Addr:    opts.ServerAddr,
		Handler: srv.Router(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			opts.Logger.Error("failed to start server", "err", err)
		}
	}()

	return &Node{
		server: server,
	}, nil
}

func (n *Node) Close() error {
	return n.server.Close()
}
