package searcher

import (
	"crypto/ecdsa"

	"github.com/lthibault/log"
	"github.com/primev/builder-boost/pkg/rollup"
)

type Config struct {
	Log            log.Logger
	Key            *ecdsa.PrivateKey
	Addr           string
	MetricsEnabled bool
	Ru             rollup.Rollup
}
