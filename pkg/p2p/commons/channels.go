package commons

type ComChannels byte

const (
	UNDEFINED_cc ComChannels = iota
	// builder channels
	SignatureCh
	BlockKeyCh
	BundleCh
	PreconfCh
	// searcher channels
	BidCh
)
