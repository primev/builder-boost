package commons

type Broadcast byte

const (
	UNDEFINED_br Broadcast = iota
	Publish
	Stream
	Gossip
)
