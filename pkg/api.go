package boost

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/attestantio/go-builder-client/api/capella"

	"github.com/lthibault/log"
	"github.com/primev/builder-boost/pkg/rollup"
	"github.com/primev/builder-boost/pkg/utils"
)

// Router paths
const (
	// proposer endpoints
	PathStatus      = "/primev/v0/status"
	PathSubmitBlock = "/primev/v1/builder/blocks"

	// searcher endpoints
	PathSearcherConnect = "/ws"
)

var (
	ErrParamNotFound = errors.New("not found")
)

type jsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type API struct {
	Service      BoostService
	Worker       *Worker
	Rollup       rollup.Rollup
	Log          log.Logger
	once         sync.Once
	mux          http.Handler
	BuilderToken string
}

func (a *API) init() {
	a.once.Do(func() {
		if a.Log == nil {
			a.Log = log.New()
		}

		// router := chi.NewRouter()
		// router.Use(middleware.Logger)
		// TODO(@floodcode): Add CORS middleware
		router := http.NewServeMux()

		router.Handle("/health", a.authenticateBuilder(http.HandlerFunc(a.handleHealthCheck)))
		// Adds an endpoint to retrieve the builder ID
		router.Handle("/builder", a.authenticateBuilder(http.HandlerFunc(a.handleBuilderID)))

		// Adds an endpoint to get commitment to the builder by searcher address
		router.HandleFunc("/commitment", a.handleSearcherCommitment)

		// TODO(@ckartik): Guard this to only by a requset made form an authorized internal service
		router.HandleFunc(PathSubmitBlock, handler(a.submitBlock))

		router.HandleFunc(PathSearcherConnect, a.ConnectedSearcher)

		a.mux = router
	})
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.init()
	a.mux.ServeHTTP(w, r)
}

func (a *API) authenticateBuilder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("X-BUIlDER-TOKEN")
		if authToken != a.BuilderToken {
			a.Log.Error("failed to authenticate builder request")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handler(f func(http.ResponseWriter, *http.Request) (int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := f(w, r)

		if status == 0 {
			status = http.StatusOK
		}

		// NOTE:  will default to http.StatusOK if f wrote any data to the
		//        response body.
		w.WriteHeader(status)

		if err != nil {
			_ = json.NewEncoder(w).Encode(jsonError{
				Code:    status,
				Message: err.Error(),
			})
		}
	}
}

type IDResponse struct {
	ID string `json:"id"`
}

// handleBuilderID returns the builder ID as an IDResponse
func (a *API) handleBuilderID(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(IDResponse{ID: a.Rollup.GetBuilderAddress().Hex()})
}

type CommitmentResponse struct {
	Commitment string `json:"commitment"`
}

func (a *API) handleSearcherCommitment(w http.ResponseWriter, r *http.Request) {
	// TODO(@ckartik): Move to middleware
	token := r.URL.Query().Get("token")
	if token == "" {
		a.Log.WithField("token", token).Error("token parameter is missing")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("token parameter is missing"))
		return
	}
	builderAddress := a.Rollup.GetBuilderAddress()

	searcherAddress, ok := utils.VerifyToken(token, builderAddress.Hex())
	if !ok {
		a.Log.WithField("token", token).Error("token is not valid")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("token is not valid"))
		return
	}

	commitment := a.Rollup.GetCommitment(searcherAddress)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(CommitmentResponse{Commitment: commitment.Hex()})
}

// connectSearcher is the handler to connect a searcher to the builder for the websocket execution hints
//
// GET /ws?token=abcd where "abcd" is the authentication token of the searcher
// The handler authenticates based on the following criteria:
// 1. The token is valid
// 2. The searcher behind the token has active subscription
// 3. The searcher behind the token is not already connected
func (a *API) ConnectedSearcher(w http.ResponseWriter, r *http.Request) {
	a.Log.Info("searcher called")
	ws := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	// Use verification scheme on token
	token := r.URL.Query().Get("token")
	if token == "" {
		a.Log.WithField("token", token).Error("token is not valid")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("token is not valid"))
		return
	}
	builderAddress := a.Rollup.GetBuilderAddress()

	searcherAddress, ok := utils.VerifyToken(token, builderAddress.Hex())
	if !ok {
		a.Log.WithField("token", token).Error("token is not valid")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("token is not valid"))
		return
	}

	_, err := a.Rollup.GetMinimalStake(builderAddress)
	if err != nil {
		if errors.Is(rollup.ErrNoMinimalStakeSet, err) {
			a.Log.WithError(err).WithField("builder_address", builderAddress).Error("no minimal stake is set, in order to allow searchers to connect, set minimal stake in the rollup contract")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		a.Log.WithError(err).Error("failed to get minimal stake")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	commitment := a.Rollup.GetCommitment(searcherAddress)
	blockNumber, err := a.Rollup.GetBlockNumber()
	if err != nil {
		a.Log.WithError(err).Error("failed to get block number")
		w.WriteHeader(http.StatusInternalServerError)
	}

	subscriptionEnd, err := a.Rollup.GetSubscriptionEnd(commitment)
	if err != nil {
		a.Log.WithError(err).Error("failed to get subscription end")
		w.WriteHeader(http.StatusInternalServerError)
	}

	searcherAddressParam := searcherAddress.Hex()
	a.Log.WithFields(logrus.Fields{"searcher": searcherAddressParam, "block_number": blockNumber, "subscription_end": subscriptionEnd}).
		Info("searcher attempting connection")

	// Check is subscription is expired
	if subscriptionEnd.Cmp(blockNumber) < 0 {
		a.Log.WithField("searcher", searcherAddressParam).
			Warn("subscription is expired")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check if searcher is already connected
	// TODO(@ckartik): Ensure we delete the searcher from the connectedSearchers map when the connection is closed
	a.Worker.lock.RLock()
	_, ok = a.Worker.connectedSearchers[searcherAddressParam]
	a.Worker.lock.RUnlock()
	if ok {
		a.Log.WithFields(logrus.Fields{"searcher": searcherAddressParam}).Error("searcher is already connected")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("searcher is already connected"))
		return
	}

	// Upgrade the HTTP request to a WebSocket connection
	conn, err := ws.Upgrade(w, r, nil)
	a.Log.Info("searcher upgraded connection")
	if err != nil {
		a.Log.Error(err)
		return
	}

	searcherConsumeChannel := make(chan Metadata, 100)
	a.Worker.lock.Lock()
	a.Worker.connectedSearchers[searcherAddressParam] = searcherConsumeChannel
	a.Worker.lock.Unlock()

	a.Log.Info("searcher connected and ready to consume data")

	closeSignalChannel := make(chan struct{})
	go func(closeChannel chan struct{}, conn *websocket.Conn) {
		for {
			a.Log.WithFields(logrus.Fields{"searcher": searcherAddressParam}).Info("starting to read from searcher")

			_, _, err := conn.NextReader()
			if err != nil {
				a.Log.WithFields(logrus.Fields{"searcher": searcherAddressParam, "err": err}).Error("error reading from searcher")
				break
			}
		}
		a.Log.WithFields(logrus.Fields{"searcher": searcherAddressParam}).Info("searcher disconnected")
		closeChannel <- struct{}{}
	}(closeSignalChannel, conn)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				a.Log.Error("recovered in searcher communication goroutine: closing connection", r)
				a.Worker.lock.Lock()
				defer a.Worker.lock.Unlock()
				delete(a.Worker.connectedSearchers, searcherAddressParam)
				conn.Close()
			}
		}()

		for {
			select {
			case <-closeSignalChannel:
				a.Worker.lock.Lock()
				defer a.Worker.lock.Unlock()
				delete(a.Worker.connectedSearchers, searcherAddressParam)
				return
			case data := <-searcherConsumeChannel:
				json, err := json.Marshal(data)
				if err != nil {
					a.Log.Error(err)
					panic(err)
				}
				conn.WriteMessage(websocket.TextMessage, json)
			}
		}
	}()

	a.Log.
		WithField("searcher_count", len(a.Worker.connectedSearchers)).
		WithField("searcher_address", searcherAddressParam).
		Info("new searcher connected")
}

// builder related handlers
func (a *API) submitBlock(w http.ResponseWriter, r *http.Request) (int, error) {
	var br capella.SubmitBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&br); err != nil {
		return http.StatusBadRequest, err
	}

	if err := a.Service.SubmitBlock(r.Context(), &br); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

type healthCheck struct {
	Searchers       []string  `json:"connected_searchers"`
	WorkerHeartBeat time.Time `json:"worker_heartbeat"`
}

// healthCheck detremines if the service is healthy
// how many connections are open
func (a *API) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// Get list of all connected searchers from a.Worker.connectedSearchers
	searchers := make([]string, 0)
	a.Worker.lock.RLock()
	for searcher := range a.Worker.connectedSearchers {
		searchers = append(searchers, searcher)
	}
	a.Worker.lock.RUnlock()

	// Send the list over the API
	json.NewEncoder(w).Encode(healthCheck{Searchers: searchers, WorkerHeartBeat: time.Unix(a.Worker.GetHeartbeat(), 0)})
}
