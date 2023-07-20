package commons

type Broadcast byte

const (
	UNDEFINED Broadcast = iota
	PUBLISH
	STREAM
	GOSSIP
)
