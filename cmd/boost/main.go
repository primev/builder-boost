package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"

	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lthibault/log"
	boost "github.com/primev/builder-boost/pkg"
	"github.com/primev/builder-boost/pkg/boostcli"
	"github.com/primev/builder-boost/pkg/rollup"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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
		Name:    "rollupkey",
		Usage:   "Private key to interact with rollup",
		Value:   "",
		EnvVars: []string{"ROLLUP_KEY"},
	},
	&cli.StringFlag{
		Name:    "rollupaddr",
		Usage:   "Rollup RPC address",
		Value:   "https://rpc.sepolia.org",
		EnvVars: []string{"ROLLUP_ADDR"},
	},
	&cli.StringFlag{
		Name:    "rollupcontract",
		Usage:   "Rollup contract address",
		Value:   "0x6e100446995f4456773Cd3e96FA201266c44d4B8",
		EnvVars: []string{"ROLLUP_CONTRACT"},
	},
	&cli.StringFlag{
		Name:    "buildertoken",
		Usage:   "Token used to authenticate request as originating from builder",
		Value:   "tmptoken",
		EnvVars: []string{"BUILDER_AUTH_TOKEN"},
	},
}

var (
	config = boost.Config{Log: log.New()}
	svr    http.Server
)

// TODO(@ckartik): Intercept SIGINT and SIGTERM to gracefully shutdown the server and persist state to cache

// Main starts the primev protocol
func main() {
	app := &cli.App{
		Name:    "builder boost",
		Usage:   "entry point to primev protocol",
		Version: version,
		Flags:   flags,
		Before:  setup(),
		Action:  run(),
	}

	if err := app.Run(os.Args); err != nil {
		config.Log.Fatal(err)
	}
}

func setup() cli.BeforeFunc {
	return func(c *cli.Context) (err error) {

		config = boost.Config{
			Log: boostcli.Logger(c),
		}

		svr = http.Server{
			Addr:           c.String("addr"),
			ReadTimeout:    c.Duration("timeout"),
			WriteTimeout:   c.Duration("timeout"),
			IdleTimeout:    time.Second * 2,
			MaxHeaderBytes: 4096,
		}

		return
	}
}

func run() cli.ActionFunc {
	return func(c *cli.Context) error {
		g, ctx := errgroup.WithContext(c.Context)

		// setup rollup service
		builderKeyString := c.String("rollupkey")
		if builderKeyString == "" {
			return errors.New("rollup key is not set, use --rollupkey option or ROLLUP_KEY env variable")
		}

		builderKeyBytes := common.FromHex(builderKeyString)
		builderKey := crypto.ToECDSAUnsafe(builderKeyBytes)

		client, err := ethclient.Dial(c.String("rollupaddr"))
		if err != nil {
			return err
		}

		contractAddress := common.HexToAddress(c.String("rollupcontract"))
		ru, err := rollup.New(client, contractAddress, builderKey, config.Log)
		if err != nil {
			return err
		}

		g.Go(func() error {
			return ru.Run(ctx)
		})

		// setup the boost service
		service := &boost.DefaultService{
			Log:    config.Log,
			Config: config,
		}

		g.Go(func() error {
			return service.Run(ctx)
		})

		// wait for the boost service to be ready
		select {
		case <-service.Ready():
		case <-ctx.Done():
			return g.Wait()
		}

		config.Log.Info("boost service ready")

		masterWorker := boost.NewWorker(service.Boost.GetWorkChannel(), config.Log)

		g.Go(func() error {
			return masterWorker.Run(ctx)
		})
		// wait for the boost service to be ready
		select {
		case <-masterWorker.Ready():
		case <-ctx.Done():
			return g.Wait()
		}

		config.Log.Info("master worker ready")

		// run the http server
		g.Go(func() (err error) {
			// set up datadog tracer
			agentAddr := c.String("agentaddr")
			if len(agentAddr) > 0 {
				tracer.Start(
					tracer.WithService("builder-boost"),
					tracer.WithEnv(c.String("env")),
					tracer.WithAgentAddr(agentAddr),
				)
				defer tracer.Stop()
			}

			svr.BaseContext = func(l net.Listener) context.Context {
				return ctx
			}

			svr.Handler = &boost.API{
				Service:      service,
				Log:          config.Log,
				Worker:       masterWorker,
				Rollup:       ru,
				BuilderToken: c.String("buildertoken"),
			}

			config.Log.Info("http server listening")
			if err = svr.ListenAndServe(); err == http.ErrServerClosed {
				err = nil
			}

			return err
		})

		g.Go(func() error {
			defer svr.Close()
			<-ctx.Done()

			ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()

			return svr.Shutdown(ctx)
		})

		return g.Wait()
	}
}
