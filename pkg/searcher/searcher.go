package searcher

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lthibault/log"

	boost "github.com/primev/builder-boost/pkg"
	"github.com/primev/builder-boost/pkg/utils"
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
	parsedURL, err := url.Parse(s.addr)
	if err != nil {
		return err
	}

	// Continue attempts to connect to the builder boost service
	token := s.GenerateAuthenticationTokenForBuilder(parsedURL)
	s.log.WithField("token", token).WithField("builder", s.addr).Info("generated token for builder boost")

	wsScheme := "ws"
	if parsedURL.Scheme == "https" {
		wsScheme = "wss"
	}

	u := url.URL{Scheme: wsScheme, Host: parsedURL.Host, Path: "/ws"}
	q := u.Query()
	q.Set("token", token)
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

type BuilderInfoResponse struct {
	BuilderID string `json:"id"`
}

// GenerateAuthenticationTokenForBuilder generates a token for the builder boost service represented by the url
func (s *searcher) GenerateAuthenticationTokenForBuilder(u *url.URL) (token string) {
	getBuilderURL := url.URL{Scheme: u.Scheme, Host: u.Host, Path: "/builder"}
	s.log.WithField("url", getBuilderURL.String()).Info("connecting to builder boost to get builder ID")

	// Make a get request to the url
	req, err := http.NewRequest("GET", getBuilderURL.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set the request header
	req.Header.Set("Content-Type", "application/json")
	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// Close the response body
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var builderInfo BuilderInfoResponse
	err = json.Unmarshal(body, &builderInfo)
	if err != nil {
		log.Fatal(err)
	}

	token, err = utils.GenerateAuthenticationToken(builderInfo.BuilderID, s.key)
	if err != nil {
		log.Fatal(err)
	}

	return token
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
		var m boost.Metadata
		err = json.Unmarshal(message, &m)
		if err != nil {
			s.log.WithField("message", string(message)).Error("failed to unmarshal message")
			continue
		}
		s.log.WithField("message", string(message)).Info("received message from builder")
	}
}
