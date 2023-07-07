package searcher

import (
	"crypto/ecdsa"

	"github.com/lthibault/log"
)

type Config struct {
	Log            log.Logger
	Key            *ecdsa.PrivateKey
	Addr           string
	MetricsEnabled bool
}
