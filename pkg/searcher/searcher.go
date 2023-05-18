package searcher

import (
	"context"
	"crypto/ecdsa"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/websocket"
	"github.com/lthibault/log"
)

type Searcher interface {
	Run(ctx context.Context) error
}

func New(config Config) Searcher {
	return &searcher{
		log:  config.Log,
		key:  config.Key,
		addr: config.Addr,
	}
}

type searcher struct {
	log  log.Logger
	key  *ecdsa.PrivateKey
	addr string
}

func (s *searcher) Run(ctx context.Context) error {
	u := url.URL{Scheme: "ws", Host: s.addr, Path: "/ws"}
	q := u.Query()
	q.Set("searcherAddress", crypto.PubkeyToAddress(s.key.PublicKey).Hex())

	u.RawQuery = q.Encode()

	for {
		s.log.WithField("url", u.String()).Info("connecting to builder boost")
		err := s.run(u.String())
		if err != nil {
			s.log.Error(err.Error())
		}

		s.log.Info("lost connection, retrying in 5 seconds")
		time.Sleep(time.Second * 5)
	}
}

func (s *searcher) run(url string) error {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	defer c.Close()

	s.log.WithField("url", url).Info("connected, processing messages")

	return s.processMessages(c)
}

func (s *searcher) processMessages(c *websocket.Conn) error {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return err
		}

		s.log.WithField("message", string(message)).Info("received message from builder")
	}
}
