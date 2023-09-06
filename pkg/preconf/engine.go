package preconf

import "log"

// Structure for preconf engine
type PreconfEngine struct {
	// Log is the logger for the preconf engine
	Log log.Logger
	// p2p *p2p.P2P Have some sort of p2p engine
}
