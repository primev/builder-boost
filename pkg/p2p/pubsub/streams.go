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
func (pss *Server) peerStreamHandler(stream network.Stream) {
	peerID := stream.Conn().RemotePeer()

	reader := bufio.NewReader(stream)
	buf := make([]byte, maxMessageSize)

	n, err := reader.Read(buf)
	if err != nil {
		return
	}

	inMsg, err := pss.imb.Parse(peerID, buf[:n])
	if err != nil {
		pss.log.With(log.F{
			"service":  "p2p stream",
			"err time": commons.GetNow(),
		}).Error(err)
	}

	// if not in approved peers try to authenticate
	if !pss.apm.InPeers(peerID) {
		pss.log.With(log.F{
			"service":  "p2p stream",
			"op":       inMsg.Op(),
			"peer":     inMsg.Peer(),
			"msg time": commons.GetNow(),
		}).Info("unverified peer message")

		switch inMsg.Op() {
		case message.Approve:
			pss.optApprove(inMsg.Peer(), inMsg.Bytes(), false)
		default:
			pss.log.With(log.F{
				"service":  "p2p stream",
				"op":       inMsg.Op(),
				"peer":     inMsg.Peer(),
				"msg time": commons.GetNow(),
			}).Info("unknown approve option!")
		}
		return
	}

	pss.log.With(log.F{
		"service":  "p2p stream",
		"op":       inMsg.Op(),
		"peer":     inMsg.Peer(),
		"msg time": commons.GetNow(),
	}).Info("verified peer message")

	switch inMsg.Op() {
	// pass auth option in this side for now
	case message.Approve:
	// create pong message and publish to show you're alive
	case message.Ping:
		pss.optPing(inMsg.Peer(), inMsg.Bytes())
	// it can be use validate peer is alive
	case message.Pong:
		pss.optPong(inMsg.Peer(), inMsg.Bytes())
	// publish version
	case message.GetVersion:
		pss.optGetVersion(inMsg.Peer())
	// store peers version
	case message.Version:
		pss.optVersion(inMsg.Peer(), inMsg.Bytes())
	// publish approved peers
	case message.GetPeerList:
		pss.optGetPeerList(inMsg.Peer())
	// it can be use for connect to peers
	case message.PeerList:
		pss.optPeerList(inMsg.Peer(), inMsg.Bytes())

	case message.PreconfBid:
		pss.optPreconfBid(inMsg.Peer(), inMsg.Bytes())
	default:
		pss.log.With(log.F{
			"service":  "p2p stream",
			"op":       inMsg.Op(),
			"peer":     inMsg.Peer(),
			"msg time": commons.GetNow(),
		}).Info("unknown option!")
	}
}
