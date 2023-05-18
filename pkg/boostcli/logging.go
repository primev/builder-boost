package boostcli

import (
	"github.com/lthibault/log"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Logger(c *cli.Context) log.Logger {
	return log.New(
		WithLevel(c),
		WithFormat(c),
		WithErrWriter(c))
}

func WithLevel(c *cli.Context) (opt log.Option) {
	var level = log.FatalLevel
	defer func() {
		opt = log.WithLevel(level)
	}()

	if c.Bool("trace") {
		level = log.TraceLevel
		return
	}

	if c.String("logfmt") == "none" {
		return
	}

	switch c.String("loglvl") {
	case "trace", "t":
		level = log.TraceLevel
	case "debug", "d":
		level = log.DebugLevel
	case "info", "i":
		level = log.InfoLevel
	case "warn", "warning", "w":
		level = log.WarnLevel
	case "error", "err", "e":
		level = log.ErrorLevel
	case "fatal", "f":
		level = log.FatalLevel
	default:
		level = log.InfoLevel
	}

	return
}

func WithFormat(c *cli.Context) log.Option {
	var fmt logrus.Formatter

	switch c.String("logfmt") {
	case "none":
	case "json":
		fmt = &logrus.JSONFormatter{PrettyPrint: c.Bool("prettyprint")}
	default:
		fmt = new(logrus.TextFormatter)
	}

	return log.WithFormatter(fmt)
}

func WithErrWriter(c *cli.Context) log.Option {
	return log.WithWriter(c.App.ErrWriter)
}
