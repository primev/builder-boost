package pubsub

import (
	"bufio"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/lthibault/log"
	"github.com/primev/builder-boost/pkg/p2p/commons"
	"github.com/primev/builder-boost/pkg/p2p/message"
)

const (
	maxMessageSize = 4096 // be careful about message size
)

// when communication needs to be between only two peers
// this streaming protocol comes into play
func (bs *BuilderServer) peerStreamHandlerB2B(stream network.Stream) {
	peerID := stream.Conn().RemotePeer()

	reader := bufio.NewReader(stream)
	buf := make([]byte, maxMessageSize)

	n, err := reader.Read(buf)
	if err != nil {
		return
	}

	inMsg, err := bs.imb.Parse(peerID, buf[:n])
	if err != nil {
		bs.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p stream b2b",
		}).Error(err)
	}

	// if not in approved peers try to authenticate
	if !bs.apm.InBuilderPeers(peerID) {
		bs.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p stream b2b",
			"op":      inMsg.Op(),
			"peer":    inMsg.Peer(),
		}).Debug("unverified peer message")

		switch inMsg.Op() {
		case message.Approve:
			bs.optApprove(bs.Stream, inMsg.Peer(), inMsg.Bytes(), false)
		default:
			bs.log.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p stream b2b",
				"op":      inMsg.Op(),
				"peer":    inMsg.Peer(),
			}).Warn("unknown approve option!")
		}
		return
	}

	bs.log.With(log.F{
		"caller":  commons.GetCallerName(),
		"date":    commons.GetNow(),
		"service": "p2p stream b2b",
		"op":      inMsg.Op(),
		"peer":    inMsg.Peer(),
	}).Debug("verified peer message")

	switch inMsg.Op() {
	// pass auth option in this side for now
	case message.Approve:
	// create pong message and publish to show you're alive
	case message.Ping:
		bs.optPing(bs.Stream, inMsg.Peer(), inMsg.Bytes())
	// it can be use validate peer is alive
	case message.Pong:
		bs.optPong(inMsg.Peer(), inMsg.Bytes())
	// publish version
	case message.GetVersion:
		bs.optGetVersion(bs.Stream, inMsg.Peer())
	// store peers version
	case message.Version:
		bs.optVersion(inMsg.Peer(), inMsg.Bytes())
	// publish approved peers
	case message.GetPeerList:
		bs.optGetPeerList(bs.Stream, inMsg.Peer())
	// it can be use for connect to peers
	case message.PeerList:
		bs.optPeerList(inMsg.Peer(), inMsg.Bytes())
	// handle the incoming signatures
	case message.Signature:
		bs.optSignature(inMsg.Peer(), inMsg.Bytes())
	// handle the incoming block keys
	case message.BlockKey:
		bs.optBlockKey(inMsg.Peer(), inMsg.Bytes())
	// handle the incoming encrypted transactions
	case message.Bundle:
		bs.optBundle(inMsg.Peer(), inMsg.Bytes())
	// handle the incoming preconf bids
	case message.PreconfBid:
		bs.optPreconfBid(inMsg.Peer(), inMsg.Bytes())

	default:
		bs.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p stream b2b",
			"op":      inMsg.Op(),
			"peer":    inMsg.Peer(),
		}).Warn("unknown option!")
	}
}

// when communication needs to be between only two peers
// this streaming protocol comes into play
// builders are not allowed to communicate with each other over this stream
func (ss *SearcherServer) peerStreamHandlerB2S(stream network.Stream) {
	peerID := stream.Conn().RemotePeer()

	reader := bufio.NewReader(stream)
	buf := make([]byte, maxMessageSize)

	n, err := reader.Read(buf)
	if err != nil {
		return
	}

	inMsg, err := ss.imb.Parse(peerID, buf[:n])
	if err != nil {
		ss.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p stream b2s",
		}).Error(err)
	}

	// if not in approved peers try to authenticate
	if ((ss.selfType == commons.Builder) && !(ss.apm.InSearcherPeers(peerID))) ||
		((ss.selfType == commons.Searcher) && !(ss.apm.InBuilderPeers(peerID))) {
		//if !ss.apm.InBuilderPeers(peerID) {
		ss.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p stream b2s",
			"op":      inMsg.Op(),
			"peer":    inMsg.Peer(),
		}).Debug("unverified peer message")

		switch inMsg.Op() {
		case message.Approve:
			ss.optApprove(ss.Stream, inMsg.Peer(), inMsg.Bytes(), false)
		default:
			ss.log.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p stream b2s",
				"op":      inMsg.Op(),
				"peer":    inMsg.Peer(),
			}).Warn("unknown approve option!")
		}
		return
	}

	ss.log.With(log.F{
		"caller":  commons.GetCallerName(),
		"date":    commons.GetNow(),
		"service": "p2p stream b2s",
		"op":      inMsg.Op(),
		"peer":    inMsg.Peer(),
	}).Debug("verified peer message")

	switch inMsg.Op() {
	// pass auth option in this side for now
	case message.Approve:
	// create pong message and publish to show you're alive
	case message.Ping:
		ss.optPing(ss.Stream, inMsg.Peer(), inMsg.Bytes())
	// it can be use validate peer is alive
	case message.Pong:
		ss.optPong(inMsg.Peer(), inMsg.Bytes())
	// publish version
	case message.GetVersion:
		ss.optGetVersion(ss.Stream, inMsg.Peer())
	// store peers version
	case message.Version:
		ss.optVersion(inMsg.Peer(), inMsg.Bytes())
	// publish approved peers
	case message.GetPeerList:
		ss.optGetPeerList(ss.Stream, inMsg.Peer())
	// it can be use for connect to peers
	case message.PeerList:
		ss.optPeerList(inMsg.Peer(), inMsg.Bytes())
	// searcher<>builder communication test
	case message.Bid:
		ss.optBid(inMsg.Peer(), inMsg.Bytes())

	default:
		ss.log.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p stream b2s",
			"op":      inMsg.Op(),
			"peer":    inMsg.Peer(),
		}).Warn("unknown option!")
	}
}
