package searcher

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lthibault/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	boost "github.com/primev/builder-boost/pkg"
	"github.com/primev/builder-boost/pkg/p2p/node"
	"github.com/primev/builder-boost/pkg/preconf"
	"github.com/primev/builder-boost/pkg/rollup"
	"github.com/primev/builder-boost/pkg/utils"
)

type Searcher interface {
	Run(ctx context.Context) error
	API(ctx context.Context) error
}

func New(config Config) Searcher {
	return &searcher{
		log:            config.Log,
		key:            config.Key,
		addr:           config.Addr,
		metricsEnabled: config.MetricsEnabled,
		ru:             config.Ru,
	}
}

type searcher struct {
	log            log.Logger
	key            *ecdsa.PrivateKey
	addr           string
	m              *metrics
	metricsEnabled bool
	ru             rollup.Rollup
	p2pEngine      node.ISearcherNode
}

type metrics struct {
	Duration prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		Duration: *prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "searcher",
			Name:      "duration",
			Help:      "Duration of the request",
			Buckets:   []float64{0.0001, 0.0002, 0.0005, 0.001, 0.002, 0.005, 0.01, 0.02, 0.03, 0.05, 0.1, 0.15, 0.2, 0.5},
		}, []string{"processing_type", "searcher_address"}),
	}

	reg.MustRegister(m.Duration)

	return m
}

// Struct for processing bid amount, ID hash, and block number
type PreconfRequest struct {
	BidAmount uint64 `json:"bid_amount"`
	IDHash    string `json:"id_hash"`
	BlockNum  uint64 `json:"block_num"`
}

func (s *searcher) preconfHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Read input from request: Bid Amount, ID Hash, Block number
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.log.WithError(err).Error("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var pc PreconfRequest
	err = json.Unmarshal(body, &pc)
	if err != nil {
		s.log.WithError(err).Error("failed to unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Construct bid
	bid, err := preconf.ConstructSignedBid(big.NewInt(int64(pc.BidAmount)), pc.IDHash, big.NewInt(int64(pc.BlockNum)), s.key)
	if err != nil {
		s.log.WithError(err).Error("failed to construct bid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bid.SubmitBid(s.p2pEngine)

	w.WriteHeader(http.StatusOK)
}

func (s *searcher) API(ctx context.Context) error {
	s.log.Info("starting searcher api")
	http.Handle("/preconf", http.HandlerFunc(s.preconfHandler))

	searchernode := node.NewSearcherNode(s.log, s.key, s.ru, nil)
	select {
	case <-searchernode.Ready():
	}
	s.p2pEngine = searchernode
	s.log.Info("searcher node ready for api")
	http.ListenAndServe(":8082", nil)

	return nil
}

func (s *searcher) Run(ctx context.Context) error {
	if s.metricsEnabled {
		reg := prometheus.NewRegistry()
		s.m = NewMetrics(reg)

		go func() {
			// Start an http server
			promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

			http.Handle("/metrics", promHandler)

			http.ListenAndServe(":8080", nil)
		}()
	}

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
	var now time.Time
	for {
		mtype, message, err := c.ReadMessage()
		if mtype == websocket.CloseMessage {
			s.log.Info("received close message from builder")
			return nil
		}
		now = time.Now()
		if err != nil {
			return err
		}
		var m boost.Metadata
		err = json.Unmarshal(message, &m)
		if err != nil {
			s.log.WithField("message", string(message)).Error("failed to unmarshal message")
			continue
		}
		if s.metricsEnabled {

			wireTime := now.Sub(m.SentTimestamp)
			e2eTime := now.Sub(m.RecTimestamp)

			s.m.Duration.WithLabelValues("wire", s.addr).Observe(wireTime.Seconds())
			s.m.Duration.WithLabelValues("e2e", s.addr).Observe(e2eTime.Seconds())

		}
		s.log.WithField("message", string(message)).Info("received message from builder")
	}
}
