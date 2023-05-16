package integrationtest

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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
	mockRollup := &rollup.MockRollup{}
	api := &boost.API{
		Rollup: mockRollup,
		Worker: boost.NewWorker(service.GetWorkChannel(), config.Log),
		Log:    config.Log,
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
	t.Run("Invalid SearcherID", func(t *testing.T) {
		invalidSearcherID := "invalidID"
		conn, resp, _ := dialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"?Searcher="+invalidSearcherID, nil)
		assert.Nil(t, conn)
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Valid SearcherID", func(t *testing.T) {
		// Setup the mock rollup
		validSearcherID := "0x812fC9524961d0566B3207fee1a567fef23E5E38"
		mockRollup.On("CheckBalance", common.HexToAddress(validSearcherID)).Return(big.NewInt(100), nil)
		mockRollup.On("GetBuilderID").Return(common.HexToAddress("0xbuilder"))
		mockRollup.On("GetMinimalStake", common.HexToAddress("0xbuilder")).Return(big.NewInt(100))

		conn, resp, _ := dialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"?Searcher="+validSearcherID, nil)
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		assert.NotNil(t, conn)
		assert.NotNil(t, resp)
	})

	t.Run("Valid Searcher ID with insufficient balance", func(t *testing.T) {
		// Setup the mock rollup
		validSearcherID := "0x812fC9524961d0566B3207fee1a567fef23E5E38"
		mockRollup.On("CheckBalance", common.HexToAddress(validSearcherID)).Return(big.NewInt(100), nil)
		mockRollup.On("GetBuilderID").Return(common.HexToAddress("0xbuilder"))
		mockRollup.On("GetMinimalStake", common.HexToAddress("0xbuilder")).Return(big.NewInt(101))

		conn, resp, _ := dialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"?Searcher="+validSearcherID, nil)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		assert.Nil(t, conn)
		assert.NotNil(t, resp)
	})

	// Test with a searcher that is already connected
	t.Run("Already connected searcher is forbidden", func(t *testing.T) {
		validSearcherID := "0x812fC9524961d0566B3207fee1a567fef23E5E38"
		mockRollup.On("CheckBalance", common.HexToAddress(validSearcherID)).Return(big.NewInt(100), nil)
		mockRollup.On("GetBuilderID").Return(common.HexToAddress("0xbuilder"))
		mockRollup.On("GetMinimalStake", common.HexToAddress("0xbuilder")).Return(big.NewInt(100))

		conn, resp, err := dialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"?Searcher="+validSearcherID, nil)
		assert.NotNil(t, conn)
		assert.NotNil(t, resp)
		assert.Nil(t, err)

		conn2, resp2, err2 := dialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"?Searcher="+validSearcherID, nil)
		assert.Equal(t, http.StatusForbidden, resp2.StatusCode)
		assert.Nil(t, conn2)
		assert.NotNil(t, resp2)
		assert.NotNil(t, err2)

	})

}
