package integrationtest

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/attestantio/go-builder-client/api/capella"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/websocket"
	"github.com/lthibault/log"
	boost "github.com/primev/builder-boost/pkg"
	"github.com/primev/builder-boost/pkg/rollup"
	"github.com/primev/builder-boost/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestConnectSearcher(t *testing.T) {
	// Initialize the API and its dependencies

	config := boost.Config{
		Log: log.New(),
	}

	// setup the boost service
	bst, err := boost.NewBoost(config)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	service := &boost.DefaultService{
		Log:    config.Log,
		Config: config,
		Boost:  bst,
	}

	go service.Run(context.TODO())

	// wait for the boost service to be ready
	<-service.Ready()
	api := &boost.API{
		Service: service,
		Worker:  boost.NewWorker(service.GetWorkChannel(), config.Log),
		Log:     config.Log,
	}

	go api.Worker.Run(context.Background())

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(api.ConnectedSearcher))
	defer server.Close()

	// WebSocket dialer
	dialer := websocket.Dialer{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}

	getWebSocketURL := func(searcherAddress common.Address) string {
		u, err := url.Parse(server.URL)
		if err != nil {
			panic(err)
		}

		q := u.Query()
		q.Set("searcherAddress", searcherAddress.Hex())

		u.Scheme = "ws"
		u.RawQuery = q.Encode()

		return u.String()
	}

	generatePrivateKey := func() (*ecdsa.PrivateKey, common.Address) {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}

		address := crypto.PubkeyToAddress(privateKey.PublicKey)

		return privateKey, address
	}

	// Test with an invalid searcher address
	t.Run("Invalid Searcher address", func(t *testing.T) {
		invalidSearcherAddress := "invalidAddress"
		conn, resp, _ := dialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"?searcherAddress="+invalidSearcherAddress, nil)
		assert.Nil(t, conn)
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Valid Searcher address", func(t *testing.T) {
		mockRollup := rollup.MockRollup{}
		_, searcherAddress := generatePrivateKey()
		builderKey, builderAddress := generatePrivateKey()
		commitment := utils.GetCommitment(builderKey, searcherAddress)

		mockRollup.On("GetBuilderAddress").Return(builderAddress)
		mockRollup.On("GetMinimalStake", builderAddress).Return(big.NewInt(100))
		mockRollup.On("GetCommitment", searcherAddress).Return(commitment)
		mockRollup.On("GetAggregaredStake", searcherAddress).Return(big.NewInt(100))
		api.Rollup = &mockRollup

		conn, resp, _ := dialer.Dial(getWebSocketURL(searcherAddress), nil)
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		assert.NotNil(t, conn)
		assert.NotNil(t, resp)
	})

	t.Run("Valid SearcherID with blocks going through to worker", func(t *testing.T) {
		// Setup the mock rollup
		mockRollup := rollup.MockRollup{}
		_, searcherAddress := generatePrivateKey()
		builderKey, builderAddress := generatePrivateKey()
		commitment := utils.GetCommitment(builderKey, searcherAddress)

		mockRollup.On("GetBuilderAddress").Return(builderAddress)
		mockRollup.On("GetMinimalStake", builderAddress).Return(big.NewInt(100))
		mockRollup.On("GetCommitment", searcherAddress).Return(commitment)
		mockRollup.On("GetAggregaredStake", searcherAddress).Return(big.NewInt(100))
		api.Rollup = &mockRollup

		conn, resp, _ := dialer.Dial(getWebSocketURL(searcherAddress), nil)
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		assert.NotNil(t, conn)
		assert.NotNil(t, resp)
		var block capella.SubmitBlockRequest
		json.Unmarshal([]byte(NoTransactionBlockRaw), &block)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func(t *testing.T) {
			// Read from connection
			mtype, r, err := conn.NextReader()
			if err != nil {
				t.Log(err)
			}
			assert.Equal(t, websocket.TextMessage, mtype)
			var data boost.Metadata
			// Decode r into data
			_ = json.NewDecoder(r).Decode(&data)
			assert.Equal(t, data.Builder, "0xaa1488eae4b06a1fff840a2b6db167afc520758dc2c8af0dfb57037954df3431b747e2f900fe8805f05d635e9a29717b")
			wg.Done()
		}(t)

		service.SubmitBlock(context.TODO(), &block)
		wg.Wait()
	})

	t.Run("Valid Searcher ID with insufficient balance", func(t *testing.T) {
		// Setup the mock rollup
		mockRollup := rollup.MockRollup{}
		_, searcherAddress := generatePrivateKey()
		builderKey, builderAddress := generatePrivateKey()
		commitment := utils.GetCommitment(builderKey, searcherAddress)

		mockRollup.On("GetBuilderAddress").Return(builderAddress)
		mockRollup.On("GetMinimalStake", builderAddress).Return(big.NewInt(101))
		mockRollup.On("GetCommitment", searcherAddress).Return(commitment)
		mockRollup.On("GetAggregaredStake", searcherAddress).Return(big.NewInt(100))
		api.Rollup = &mockRollup

		_, resp, _ := dialer.Dial(getWebSocketURL(searcherAddress), nil)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		assert.NotNil(t, resp)
	})

	// Test with a searcher that is already connected
	t.Run("Already connected searcher is forbidden", func(t *testing.T) {
		mockRollup := rollup.MockRollup{}
		_, searcherAddress := generatePrivateKey()
		builderKey, builderAddress := generatePrivateKey()
		commitment := utils.GetCommitment(builderKey, searcherAddress)

		mockRollup.On("GetBuilderAddress").Return(builderAddress)
		mockRollup.On("GetMinimalStake", builderAddress).Return(big.NewInt(100))
		mockRollup.On("GetCommitment", searcherAddress).Return(commitment)
		mockRollup.On("GetAggregaredStake", searcherAddress).Return(big.NewInt(100))
		api.Rollup = &mockRollup

		conn, resp, err := dialer.Dial(getWebSocketURL(searcherAddress), nil)
		assert.NotNil(t, conn)
		assert.NotNil(t, resp)
		assert.Nil(t, err)

		conn2, resp2, err2 := dialer.Dial(getWebSocketURL(searcherAddress), nil)
		assert.Equal(t, http.StatusForbidden, resp2.StatusCode)
		assert.Nil(t, conn2)
		assert.NotNil(t, resp2)
		assert.NotNil(t, err2)
	})

	// Test with a searcher that is already connected
	// TODO(@ckartik): Resolve this test as a searcher who tries to reconnect will fail
	t.Run("Already connected searcher is closes connection and reopens", func(t *testing.T) {
		mockRollup := rollup.MockRollup{}
		_, searcherAddress := generatePrivateKey()
		builderKey, builderAddress := generatePrivateKey()
		commitment := utils.GetCommitment(builderKey, searcherAddress)

		mockRollup.On("GetBuilderAddress").Return(builderAddress)
		mockRollup.On("GetMinimalStake", builderAddress).Return(big.NewInt(100))
		mockRollup.On("GetCommitment", searcherAddress).Return(commitment)
		mockRollup.On("GetAggregaredStake", searcherAddress).Return(big.NewInt(100))
		api.Rollup = &mockRollup

		conn, resp, err := dialer.Dial(getWebSocketURL(searcherAddress), nil)
		assert.NotNil(t, conn)
		assert.NotNil(t, resp)
		assert.Nil(t, err)
		conn.Close()

		_, resp2, err2 := dialer.Dial(getWebSocketURL(searcherAddress), nil)
		assert.Equal(t, http.StatusSwitchingProtocols, resp2.StatusCode)
		assert.NotNil(t, resp2)
		assert.Nil(t, err2)
	})
}
