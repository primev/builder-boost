package commons

type PeerType byte

const (
	UNDEFINED_pt PeerType = iota
	Builder
	Searcher
)

func (p PeerType) String() string {
	switch p {
	case Builder:
		return "builder"
	case Searcher:
		return "searcher"
	default:
		return "undefined!"
	}
}
