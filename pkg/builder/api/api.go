package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/attestantio/go-builder-client/api/capella"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
	"github.com/primev/builder-boost/pkg/apiserver"
	"github.com/primev/builder-boost/pkg/builder/searcherclient"
	"github.com/primev/builder-boost/pkg/rollup"
	"github.com/primev/builder-boost/pkg/utils"
	"golang.org/x/exp/slog"
)

const (
	defaultNamespace = "primev"
)

type searcherKey struct{}

type api struct {
	builderToken string
	rollUp       rollup.Rollup
	logger       *slog.Logger
	sclient      searcherclient.SearcherClient
}

type APIServer interface {
	ChainHandlers(string, http.Handler, ...func(http.Handler) http.Handler)
}

// RegisterAPI registers the API handlers with the provided server. It doesnt
// return anything as it is assumed that the server will be started after
// registration. Lifecycle of arguments is the responsibility of the caller.
func RegisterAPI(
	token string,
	server APIServer,
	rollUp rollup.Rollup,
	logger *slog.Logger,
	sclient searcherclient.SearcherClient,
) {

	a := &api{
		builderToken: token,
		rollUp:       rollUp,
		logger:       logger,
		sclient:      sclient,
	}
	server.ChainHandlers(
		"/health",
		apiserver.MethodHandler(http.MethodGet, a.handleHealthCheck),
		a.authSearcher,
	)

	server.ChainHandlers(
		"/builder",
		apiserver.MethodHandler(http.MethodGet, a.handleBuilderID),
		a.authenticateBuilder,
	)

	server.ChainHandlers(
		"/commitment",
		apiserver.MethodHandler(http.MethodGet, a.handleSearcherCommitment),
		a.authSearcher,
	)

	server.ChainHandlers(
		"/primev/v1/builder/blocks",
		apiserver.MethodHandler(http.MethodPost, a.submitBlock),
	)

	server.ChainHandlers(
		"/ws",
		apiserver.MethodHandler(http.MethodGet, a.connectSearcher),
	)
}

// IDResponse is a simple struct for returning an ID
type IDResponse struct {
	ID string `json:"id"`
}

func (a *api) authenticateBuilder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("X-Builder-Token")
		if authToken != a.builderToken {
			a.logger.Error("failed to authenticate builder request")
			err := apiserver.WriteResponse(w, http.StatusUnauthorized, "token invalid")
			if err != nil {
				a.logger.Error("error writing response", "err", err)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleBuilderID returns the builder ID as an IDResponse
func (a *api) handleBuilderID(w http.ResponseWriter, r *http.Request) {
	logger := a.logger.With("method", "handleBuilderID")
	resp := IDResponse{ID: a.rollUp.GetBuilderAddress().Hex()}

	err := apiserver.WriteResponse(w, http.StatusOK, resp)
	if err != nil {
		logger.Error("error writing response", "err", err)
	}
}

// CommitmentResponse is a simple struct for returning a commitment
type CommitmentResponse struct {
	Commitment string `json:"commitment"`
}

func (a *api) authSearcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("X-Primev-Signature")
		builderAddress := a.rollUp.GetBuilderAddress()
		searcherAddress, ok := utils.VerifyAuthenticationToken(authToken, builderAddress.Hex())
		if !ok {
			a.logger.Error(
				"error verifying authentication token",
				"token", authToken,
				"builderAddress", builderAddress.Hex(),
			)
			err := apiserver.WriteResponse(w, http.StatusUnauthorized, "token is not valid")
			if err != nil {
				a.logger.Error("error writing response", "err", err)
			}
			return
		}

		reqClone := r.Clone(context.WithValue(r.Context(), searcherKey{}, searcherAddress))

		next.ServeHTTP(w, reqClone)
	})
}

// handleSearcherCommitment returns the searcher commitment as a CommitmentResponse
func (a *api) handleSearcherCommitment(w http.ResponseWriter, r *http.Request) {
	logger := a.logger.With("method", "handleSearcherCommitment")

	searcherAddress, ok := r.Context().Value(searcherKey{}).(common.Address)
	if !ok {
		logger.Error("error getting searcher address from context")
		// This should never happen
		err := apiserver.WriteResponse(
			w,
			http.StatusInternalServerError,
			"searcher address not found",
		)
		if err != nil {
			logger.Error("error writing response", "err", err)
		}
		return
	}

	commitment := CommitmentResponse{Commitment: a.rollUp.GetCommitment(searcherAddress).Hex()}
	err := apiserver.WriteResponse(w, http.StatusOK, commitment)
	if err != nil {
		logger.Error("error writing response", "err", err)
	}
}

// connectSearcher is the handler to connect a searcher to the builder for the websocket execution hints
// TODO(@ckartik): Move the handling of searcher connection to service layer
//
// GET /ws?token=abcd where "abcd" is the authentication token of the searcher
// The handler authenticates based on the following criteria:
// 1. The token is valid
// 2. The searcher behind the token has active subscription
// 3. The searcher behind the token is not already connected
func (a *api) connectSearcher(w http.ResponseWriter, r *http.Request) {
	logger := a.logger.With("method", "connectSearcher")

	// Use verification scheme on token
	token := r.URL.Query().Get("token")
	if token == "" {
		logger.Error("token is not provided")
		err := apiserver.WriteResponse(w, http.StatusBadRequest, "token is not provided")
		if err != nil {
			logger.Error("error writing response", "err", err)
		}
		return
	}

	builderAddress := a.rollUp.GetBuilderAddress()

	searcherAddress, ok := utils.VerifyAuthenticationToken(token, builderAddress.Hex())
	if !ok {
		logger.Error("token is not valid", "token", token)
		err := apiserver.WriteResponse(w, http.StatusForbidden, "token is not valid")
		if err != nil {
			logger.Error("error writing response", "err", err)
		}
		return
	}

	_, err := a.rollUp.GetMinimalStake(builderAddress)
	if err != nil {
		if errors.Is(rollup.ErrNoMinimalStakeSet, err) {
			logger.Error(
				"no minimal stake in the rollup contract",
				"builder_address", builderAddress,
			)
			err := apiserver.WriteResponse(
				w,
				http.StatusForbidden,
				"no minimal stake in the rollup contract",
			)
			if err != nil {
				logger.Error("error writing response", "err", err)
			}
			return
		}
		logger.Error("failed to get minimal stake", "err", err)
		err := apiserver.WriteResponse(
			w,
			http.StatusInternalServerError,
			"failed to get minimal stake",
		)
		if err != nil {
			logger.Error("error writing response", "err", err)
		}
		return
	}

	commitment := a.rollUp.GetCommitment(searcherAddress)
	blockNumber, err := a.rollUp.GetBlockNumber()
	if err != nil {
		logger.Error("failed to get block number", "err", err)
		err := apiserver.WriteResponse(
			w,
			http.StatusInternalServerError,
			"failed to get block number",
		)
		if err != nil {
			logger.Error("error writing response", "err", err)
		}
		return
	}

	subscriptionEnd, err := a.rollUp.GetSubscriptionEnd(commitment)
	if err != nil {
		logger.Error("failed to get subscription end", "err", err)
		err := apiserver.WriteResponse(
			w,
			http.StatusInternalServerError,
			"failed to get subscription end",
		)
		if err != nil {
			logger.Error("error writing response", "err", err)
		}
		return
	}

	// Check is subscription is expired
	if subscriptionEnd.Cmp(blockNumber) < 0 {
		logger.Error("subscription is expired", "searcher", searcherAddress.Hex())
		err := apiserver.WriteResponse(
			w,
			http.StatusForbidden,
			"subscription is expired",
		)
		if err != nil {
			logger.Error("error writing response", "err", err)
		}
		return
	}

	if a.sclient.IsConnected(searcherAddress.Hex()) {
		logger.Error("searcher is already connected, closing old connection",
			"searcher", searcherAddress.Hex(),
		)
		a.sclient.RemoveSearcher(searcherAddress.Hex())
	}

	ws := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := ws.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("failed to upgrade connection", "err", err)
		err := apiserver.WriteResponse(
			w,
			http.StatusInternalServerError,
			"failed to upgrade connection",
		)
		if err != nil {
			logger.Error("error writing response", "err", err)
		}
		return
	}

	searcher := searcherclient.NewSearcher(
		searcherAddress.Hex(),
		conn,
		a.rollUp,
		a.logger.With("searcher", searcherAddress.Hex()),
		subscriptionEnd.Int64(),
		func(reason string) { a.sclient.Disconnected(searcherAddress.Hex(), reason) },
	)
	a.sclient.AddSearcher(searcher)

	logger.Info("searcher attempting connection",
		"searcher", searcherAddress.Hex(),
		"block_number", blockNumber,
		"subscription_end", subscriptionEnd,
	)
}

// builder related handlers
func (a *api) submitBlock(w http.ResponseWriter, r *http.Request) {
	logger := a.logger.With("method", "submitBlock")

	br, err := apiserver.BindJSON[capella.SubmitBlockRequest](w, r)
	if err != nil {
		logger.Error("failed to decode submit block request", "err", err)
		return
	}

	if err := a.sclient.SubmitBlock(r.Context(), &br); err != nil {
		logger.Error("failed to submit block", "err", err)
		err := apiserver.WriteResponse(
			w,
			http.StatusInternalServerError,
			"failed to submit block",
		)
		if err != nil {
			logger.Error("error writing response", "err", err)
		}
		return
	}

	err = apiserver.WriteResponse(w, http.StatusOK, "block submitted")
	if err != nil {
		logger.Error("error writing response", "err", err)
	}
}

type healthCheck struct {
	Searchers       []string  `json:"connected_searchers"`
	WorkerHeartBeat time.Time `json:"worker_heartbeat"`
}

// healthCheck detremines if the service is healthy
// how many connections are open
func (a *api) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	logger := a.logger.With("method", "handleHealthCheck")
	sInfo := a.sclient.GetSeacherInfo()

	err := apiserver.WriteResponse(w, http.StatusOK, sInfo)
	if err != nil {
		logger.Error("error writing response", "err", err)
	}
}
