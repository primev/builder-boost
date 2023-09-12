package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/primev/builder-boost/pkg/builder"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

var (
	// version is the current version of the application. It is set during build.
	version string = "dev"
	// commit is the current commit of the application. It is set during build.
	commit string = "dirty"
	// commitTime is the time the application was built. It is set during build.
	commitTime string
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
		Name:    "addr",
		Usage:   "server listen address",
		Value:   ":18550",
		EnvVars: []string{"BOOST_ADDR"},
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
		Name:     "rollupkey",
		Usage:    "Private key to interact with rollup",
		Value:    "",
		EnvVars:  []string{"ROLLUP_KEY"},
		Required: true,
	},
	&cli.StringFlag{
		Name:    "rollupaddr",
		Usage:   "Rollup RPC address",
		Value:   "https://ethereum-sepolia.blockpi.network/v1/rpc/public",
		EnvVars: []string{"ROLLUP_ADDR"},
	},
	&cli.StringFlag{
		Name:    "rollupcontract",
		Usage:   "Rollup contract address",
		Value:   "0x6219a236EFFa91567d5ba4a0A5134297a35b0b2A",
		EnvVars: []string{"ROLLUP_CONTRACT"},
	},
	&cli.StringFlag{
		Name:     "buildertoken",
		Usage:    "Token used to authenticate request as originating from builder",
		Value:    "",
		EnvVars:  []string{"BUILDER_AUTH_TOKEN"},
		Required: true,
	},
	&cli.BoolFlag{
		Name:    "metrics",
		Usage:   "enables metrics tracking for boost",
		Value:   false,
		EnvVars: []string{"METRICS"},
	},
	&cli.StringFlag{
		Name:    "dacontract",
		Usage:   "DA contract address",
		Value:   "0xac27A2cbdBA8768D49e359ebA326fC1F27832ED4",
		EnvVars: []string{"ROLLUP_CONTRACT"},
	},
	&cli.StringFlag{
		Name:    "daaddr",
		Usage:   "DA RPC address",
		Value:   "http://54.200.76.18:8545",
		EnvVars: []string{"ROLLUP_ADDR"},
	},
	&cli.BoolFlag{
		Name:    "inclusionlist",
		Usage:   "enables inclusion list for boost",
		Value:   false,
		EnvVars: []string{"INCLUSION_LIST"},
	},
}

func main() {

	// Parse commit time
	commitTimeInt, err := strconv.ParseInt(commitTime, 10, 64)
	if err != nil {
		commitTimeInt = time.Now().Unix()
	}

	app := &cli.App{
		Name:     "builder boost",
		Version:  fmt.Sprintf("%s-%s", version, commit),
		Compiled: time.Unix(commitTimeInt, 0),
		Usage:    "entry point to primev protocol",
		Flags:    flags,
		Action:   run(),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(app.Writer, "exited with error: %v\n", err)
	}
}

func run() cli.ActionFunc {
	return func(c *cli.Context) error {
		logger, err := newLogger(c.String("loglvl"), c.String("logfmt"), c.App.Writer)
		if err != nil {
			return fmt.Errorf("failed to create logger: %w", err)
		}

		nd, err := builder.NewNode(
			c.Context,
			builder.Options{
				Logger:               logger,
				ServerAddr:           c.String("addr"),
				TracingEndpoint:      c.String("agentaddr"),
				RollupKey:            c.String("rollupkey"),
				RollupAddr:           c.String("rollupaddr"),
				RollupContract:       c.String("rollupcontract"),
				BuilderToken:         c.String("buildertoken"),
				MetricsEnabled:       c.Bool("metrics"),
				DAContract:           c.String("dacontract"),
				DARPCAddr:            c.String("daaddr"),
				InclusionListEnabled: c.Bool("inclusionlist"),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create node: %w", err)
		}

		<-c.Done()
		fmt.Fprintf(c.App.Writer, "shutting down...\n")
		closed := make(chan struct{})

		go func() {
			defer close(closed)

			err := nd.Close()
			if err != nil {
				logger.Error("failed to close node", "error", err)
			}
		}()

		select {
		case <-closed:
		case <-time.After(5 * time.Second):
			logger.Error("failed to close node in time")
		}

		return nil
	}
}

func newLogger(lvl, logFmt string, sink io.Writer) (*slog.Logger, error) {
	var (
		level   = new(slog.LevelVar) // Info by default
		handler slog.Handler
	)

	switch lvl {
	case "debug":
		level.Set(slog.LevelDebug)
	case "info":
		level.Set(slog.LevelInfo)
	case "warn":
		level.Set(slog.LevelWarn)
	case "error":
		level.Set(slog.LevelError)
	default:
		return nil, fmt.Errorf("invalid log level: %s", lvl)
	}

	switch logFmt {
	case "text":
		handler = slog.NewTextHandler(sink, &slog.HandlerOptions{Level: level})
	case "none":
		fallthrough
	case "json":
		handler = slog.NewJSONHandler(sink, &slog.HandlerOptions{Level: level})
	default:
		return nil, fmt.Errorf("invalid log format: %s", logFmt)
	}

	return slog.New(handler), nil
}
