package stream

import (
	"bufio"
	"context"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// TODO research streams: Can communication continue without closing a stream?
// Is this a correct approach? If so, how?"

// stream is an interface for sending messages over a network stream
type Stream interface {
	Send(context.Context, peer.ID, []byte) error
	Close() error
}

type stream struct {
	// the libp2p host
	host host.Host

	// The protocol ID
	proto protocol.ID

	// the network stream (NOTE:unused for now)
	stream network.Stream
}

// it creates a new Stream instance and sets the StreamHandler
func New(
	host host.Host,
	proto protocol.ID,
	handler func(stream network.Stream),
) Stream {
	host.SetStreamHandler(proto, handler)

	return &stream{
		host:  host,
		proto: proto,
	}
}

// it sends a message to the specified peer over the given protocol
func (s *stream) Send(ctx context.Context, p peer.ID, msg []byte) error {
	stream, err := s.host.NewStream(ctx, p, s.proto)
	if err != nil {
		return err
	}

	defer stream.Close()

	writer := bufio.NewWriter(stream)
	_, err = writer.Write(msg)
	if err != nil {
		return err
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	return nil
}

// it closes the stream
func (s *stream) Close() error {
	if s.stream != nil {
		err := s.stream.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
