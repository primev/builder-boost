package boost

import (
	"github.com/lthibault/log"
)

// Config provides all available options for auction-gateway
type Config struct {
	Log                  log.Logger
	InclusionProofActive bool
}

// noop for now
func (c *Config) validate() error {
	return nil
}

// func ensureScheme(url string) string {
// 	if strings.HasPrefix(url, "http") {
// 		return url
// 	}

// 	return fmt.Sprintf("http://%s", url)
// }
