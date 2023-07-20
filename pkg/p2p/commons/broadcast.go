package commons

type Broadcast byte

const (
	UNDEFINED Broadcast = iota
	Publish
	Stream
	Gossip
)
