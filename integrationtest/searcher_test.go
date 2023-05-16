package integrationtest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/lthibault/log"
	boost "github.com/primev/builder-boost/pkg"
	"github.com/primev/builder-boost/pkg/rollup"
	"github.com/stretchr/testify/assert"
)

func TestConnectSearcher(t *testing.T) {
	// Initialize the API and its dependencies

	config := boost.Config{
		Log: log.New(),
	}
	// setup the boost service
	service := &boost.DefaultService{
		Log:    config.Log,
		Config: config,
	}

	go service.Run(context.TODO())

	// wait for the boost service to be ready
	<-service.Ready()

	api := &boost.API{
		Rollup: &rollup.MockRollup{},
		Worker: boost.NewWorker(service.GetWorkChannel(), config.Log),
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(api.ConnectedSearcher))
	defer server.Close()

	// WebSocket dialer
	dialer := websocket.Dialer{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}

	// Test withan invalid searcher ID
	t.Run("InvalidSearcherID", func(t *testing.T) {
		invalidSearcherID := "invalidID"
		conn, resp, err := dialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"?Searcher="+invalidSearcherID, nil)
		assert.Nil(t, conn)
		assert.NotNil(t, resp)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// // Test with a searcher having insufficient balance
	// t.Run("InsufficientBalance", func(t *testing.T) {
	// 	insufficientBalanceSearcherID := "0x1234567890123456789012345678901234567890"
	// 	api.Rollup.GetMinimalStake = func(builderID common.Address) *big.Int {
	// 		return big.NewInt(100)
	// 	}
	// 	api.Rollup.CheckBalance = func(searcherID common.Address) *big.Int {
	// 		return big.NewInt(50)
	// 	}
	// 	conn, resp, err := dialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"?Searcher="+insufficientBalanceSearcherID, nil)
	// 	assert.Nil(t, conn)
	// 	assert.NotNil(t, resp)
	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	// })

	// // Test with a searcher that is already connected
	// t.Run("AlreadyConnected", func(t *testing.T) {
	// 	alreadyConnectedSearcherID := "0x1234567890123456789012345678901234567891"
	// 	api.Worker.connectedSearchers[alreadyConnectedSearcherID] = make(chan Metadata, 100)
	// 	conn, resp, err := dialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"?Searcher="+alreadyConnectedSearcherID, nil)
	// 	assert.Nil(t, conn)
	// 	assert.NotNil(t, resp)
	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	// })

	// // Test with a valid searcher ID and sufficient balance
	// t.Run("ValidSearcherID", func(t *testing.T) {
	// 	validSearcherID := "0x1234567890123456789012345678901234567892"
	// 	api.Rollup.GetMinimalStake = func(builderID common.Address) *big.Int {
	// 		return big.NewInt(100)
	// 	}
	// 	api.Rollup.CheckBalance = func(searcherID common.Address) *big.Int {
	// 		return big.NewInt(200)
	// 	}
	// 	conn, resp, err := dialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"?Searcher="+validSearcherID, nil)
	// 	assert.NotNil(t, conn)
	// 	assert.Nil(t, resp)
	// 	assert.Nil(t, err)

	// 	// Send a message to the searcher
	// 	metadata := Metadata{
	// 		// Fill in the metadata fields
	// 	}
	// 	api.Worker.connectedSearchers[validSearcherID] <- metadata

	// 	// Read the message from the WebSocket connection
	// 	_, msg, err := conn.ReadMessage()
	// 	assert.Nil(t, err)

	// 	// Unmarshal the message and compare with the original metadata
	// 	var receivedMetadata Metadata
	// 	err = json.Unmarshal(msg, &receivedMetadata)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, metadata, receivedMetadata)

	// 	// Close the connection
	// 	conn.Close()
	// })
}
