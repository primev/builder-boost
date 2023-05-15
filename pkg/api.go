package boost

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/attestantio/go-builder-client/api/capella"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lthibault/log"
	"github.com/primev/builder-boost/pkg/rollup"
	// "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

// Router paths
const (
	// proposer endpoints
	PathStatus          = "/primev/v0/status"
	PathSubmitBlock     = "/primev/v1/builder/blocks"
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
	Service   BoostService
	Worker    *Worker
	Rollup    rollup.Rollup
	Log       log.Logger
	once      sync.Once
	mux       http.Handler
	BuilderID common.Address
}

func (a *API) init() {
	a.once.Do(func() {
		if a.Log == nil {
			a.Log = log.New()
		}

		// TODO(@floodcode): Add CORS middleware

		router := http.NewServeMux()

		// router.Use(
		// 	withDrainBody(),
		// 	withContentType("application/json"),
		// 	withLogger(a.Log),
		// ) // set middleware

		// root returns 200 - nil
		router.HandleFunc("/", succeed(http.StatusOK))

		// proposer related
		// router.HandleFunc(PathStatus, succeed(http.StatusOK)).Methods(http.MethodGet)

		// TODO(@ckartik): Guard this to only by a requset made form an authorized internal service
		// builder related
		router.HandleFunc(PathSubmitBlock, handler(a.submitBlock))

		router.HandleFunc(PathSearcherConnect, a.ConnectedSearcher)

		a.mux = router
	})

}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.init()
	a.mux.ServeHTTP(w, r)
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

// connectSearcher is the handler to connect a searcher to the builder for the websocket execution hints
// TODO: Add authentication
//
// GET /ws?Searcher=0x123 where 0x123 is the address of the searcher (soon to be auth token)
// The handler authenticates based on the following criteria:
// 1. The address is a valid address
// 2. The address has sufficient balance
// 3. The address is not already connected

func (a *API) ConnectedSearcher(w http.ResponseWriter, r *http.Request) {
	log.Info("searcher called")
	ws := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	searcherID := r.URL.Query().Get("Searcher")
	if !common.IsHexAddress(searcherID) {
		log.Error("Searcher ID is not a valid address", "searcherID", searcherID)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Searcher ID is not a valid address"))
		return
	}

	balance := a.Rollup.CheckBalance(common.HexToAddress(searcherID))
	log.Info("Searcher attempting connection", "searcherID", searcherID, "balance", balance)

	// Check for sufficent balance
	if balance.Cmp(a.Rollup.GetMinimalStake(a.Rollup.GetBuilderID())) < 0 {
		log.Error("Searcher has insufficient balance", "balance", balance, "required", a.Rollup.GetMinimalStake(a.Rollup.GetBuilderID()))
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Searcher has insufficient balance"))
		return
	}

	// Check if searcher is already connected
	// TODO(@ckartik): Ensure we delete the searcher from the connectedSearchers map when the connection is closed
	a.Worker.lock.RLock()
	_, ok := a.Worker.connectedSearchers[searcherID]
	a.Worker.lock.RUnlock()
	if ok {
		log.Error("Searcher is already connected", "searcherID", searcherID)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Searcher is already connected"))
		return
	}

	// Upgrade the HTTP request to a WebSocket connection
	conn, err := ws.Upgrade(w, r, nil)
	log.Info("searcher upgraded connection")
	if err != nil {
		log.Error(err)
		return
	}

	searcherConsumeChannel := make(chan Metadata, 100)
	a.Worker.lock.Lock()
	a.Worker.connectedSearchers[searcherID] = searcherConsumeChannel
	a.Worker.lock.Unlock()

	log.Info("Searcher connected and ready to consume data")

	// Listen for messages from the client
	// TODO(@ckartik): Turn into a select statment?
	for {
		data := <-searcherConsumeChannel
		json, err := json.Marshal(data)
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("Sending message", "msg", json)
		conn.WriteMessage(websocket.TextMessage, json)
	}

	// // Close the connection
	// conn.Close()
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

func succeed(status int) http.HandlerFunc {
	return handler(func(http.ResponseWriter, *http.Request) (int, error) {
		return status, nil
	})
}

// healthCheck detremines if the service is healthy
// how many connections are open
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
