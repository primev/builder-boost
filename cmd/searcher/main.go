package main

import (
	"errors"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lthibault/log"
	boost "github.com/primev/builder-boost/pkg"
	"github.com/primev/builder-boost/pkg/boostcli"
	"github.com/primev/builder-boost/pkg/searcher"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

const (
	shutdownTimeout = 5 * time.Second
	version         = boost.Version
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "loglvl",
		Usage:   "logging level: trace, debug, info, warn, error or fatal",
		Value:   "info",
		EnvVars: []string{"LOGLVL"},
	},
	&cli.StringFlag{
		Name:    "logfmt",
		Usage:   "format logs as text, json or none",
		Value:   "text",
		EnvVars: []string{"LOGFMT"},
	},
	&cli.StringFlag{
		Name:    "boosturl",
		Usage:   "boost endpoint url",
		Value:   "http://localhost:18550",
		EnvVars: []string{"BOOST_URL"},
	},
	&cli.StringFlag{
		Name:    "env",
		Usage:   "service environment (development, production, etc.)",
		Value:   "development",
		EnvVars: []string{"ENV"},
	},
	&cli.StringFlag{
		Name:    "agentaddr",
		Usage:   "datadog agent address",
		Value:   "",
		EnvVars: []string{"AGENT_ADDR"},
	},
	&cli.StringFlag{
		Name:    "searcherkey",
		Usage:   "private key to interact with builder boost",
		Value:   "",
		EnvVars: []string{"SEARCHER_KEY"},
	},
	&cli.BoolFlag{
		Name:    "metrics",
		Usage:   "enable metrics",
		Value:   false,
		EnvVars: []string{"METRICS"},
	},
}
var (
	config = searcher.Config{Log: log.New()}
)

func main() {
	app := &cli.App{
		Name:    "builder boost searcher",
		Usage:   "entry point to primev protocol",
		Version: version,
		Flags:   flags,
		Action:  run(),
	}

	if err := app.Run(os.Args); err != nil {
		config.Log.Fatal(err)
	}
}

func run() cli.ActionFunc {
	return func(c *cli.Context) error {
		g, ctx := errgroup.WithContext(c.Context)

		searcherKeyString := c.String("searcherkey")
		if searcherKeyString == "" {
			return errors.New("searcher key is not set, use --searcherkey option or SEARCHER_KEY env variable")
		}

		boostAddrString := c.String("boosturl")
		if boostAddrString == "" {
			return errors.New("boost URL is not set, use --boosturl option or BOOST_URL env variable")
		}

		searcherKeyBytes := common.FromHex(searcherKeyString)
		searcherKey := crypto.ToECDSAUnsafe(searcherKeyBytes)

		config.Log = boostcli.Logger(c)
		config.Key = searcherKey
		config.Addr = boostAddrString
		config.MetricsEnabled = c.Bool("metrics")
		searcher := searcher.New(config)

		g.Go(func() error {
			return searcher.Run(ctx)
		})

		return g.Wait()
	}
}
