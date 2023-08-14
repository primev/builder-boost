package commons

import "github.com/libp2p/go-libp2p/core/peer"

type ReportType byte

// TODO: additional reporting types might be necessary, consult the team
const (
	UNDEFINED_rt ReportType = iota
	BadSignatureMsg
	BadBlockKeyMsg
	BadBundleMsg
	BadPreconfMsg
)

func (r ReportType) String() string {
	switch r {
	case BadSignatureMsg:
		return "bad_signature_msg"
	case BadBlockKeyMsg:
		return "bad_block_key_msg"
	case BadBundleMsg:
		return "bad_bundle_msg"
	case BadPreconfMsg:
		return "bad_preconf_msg"
	default:
		return "unknown report type"
	}
}

// ReportEvent represents an event involving a reported issue
type ReportEvent struct {
	PeerID peer.ID
	Reason ReportType
}
