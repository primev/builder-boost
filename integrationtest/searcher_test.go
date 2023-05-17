package integrationtest

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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
	service := &boost.DefaultService{
		Log:    config.Log,
		Config: config,
	}

	go service.Run(context.TODO())

	// wait for the boost service to be ready
	<-service.Ready()
	api := &boost.API{
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

	getWebSocketURL := func(searcherAddress, commitmentAddress common.Address) string {
		u, err := url.Parse(server.URL)
		if err != nil {
			panic(err)
		}

		q := u.Query()
		q.Set("searcherAddress", searcherAddress.Hex())
		q.Set("commitmentAddress", commitmentAddress.Hex())

		u.Scheme = "ws"
		u.RawQuery = q.Encode()

		return u.String()
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
		// Setup the mock rollup
		mockRollup := rollup.MockRollup{}
		validSearcherAddress := common.HexToAddress("0x812fC9524961d0566B3207fee1a567fef23E5E38")
		validCommitmentAddress := common.HexToAddress("0x812fC9524961d0566B3207fee1a567fef23E5E38")
		builderAddress := common.HexToAddress("0xbuilder")
		commitment := utils.GetCommitment(validCommitmentAddress, builderAddress)
		mockRollup.On("GetStake", validSearcherAddress, commitment).Return(big.NewInt(100), nil)
		mockRollup.On("GetBuilderAddress").Return(builderAddress)
		mockRollup.On("GetMinimalStake", builderAddress).Return(big.NewInt(100))
		api.Rollup = &mockRollup

		conn, resp, _ := dialer.Dial(getWebSocketURL(validSearcherAddress, validCommitmentAddress), nil)
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		assert.NotNil(t, conn)
		assert.NotNil(t, resp)
	})

	t.Run("Valid Searcher address with insufficient balance", func(t *testing.T) {
		// Setup the mock rollup
		mockRollup := rollup.MockRollup{}
		validSearcherAddress := common.HexToAddress("0x812fC9524961d0566B3207fee1a567fef23E5E38")
		validCommitmentAddress := common.HexToAddress("0x812fC9524961d0566B3207fee1a567fef23E5E38")
		builderAddress := common.HexToAddress("0xbuilder2")
		commitment := utils.GetCommitment(validCommitmentAddress, builderAddress)
		mockRollup.On("GetStake", validSearcherAddress, commitment).Return(big.NewInt(100), nil)
		mockRollup.On("GetBuilderAddress").Return(builderAddress)
		mockRollup.On("GetMinimalStake", builderAddress).Return(big.NewInt(101))
		api.Rollup = &mockRollup

		_, resp, _ := dialer.Dial(getWebSocketURL(validSearcherAddress, validCommitmentAddress), nil)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		assert.NotNil(t, resp)
	})

	// Test with a searcher that is already connected
	t.Run("Already connected searcher is forbidden", func(t *testing.T) {
		mockRollup := rollup.MockRollup{}
		validSearcherAddress := common.HexToAddress("0x812fC9524961d0566B3207fee1a567fef23E5E39")
		validCommitmentAddress := common.HexToAddress("0x812fC9524961d0566B3207fee1a567fef23E5E38")
		builderAddress := common.HexToAddress("0xbuilder")
		commitment := utils.GetCommitment(validCommitmentAddress, builderAddress)

		mockRollup.On("GetStake", validSearcherAddress, commitment).Return(big.NewInt(100), nil)
		mockRollup.On("GetBuilderAddress").Return(builderAddress)
		mockRollup.On("GetMinimalStake", builderAddress).Return(big.NewInt(100))
		api.Rollup = &mockRollup

		conn, resp, err := dialer.Dial(getWebSocketURL(validSearcherAddress, validCommitmentAddress), nil)
		assert.NotNil(t, conn)
		assert.NotNil(t, resp)
		assert.Nil(t, err)

		conn2, resp2, err2 := dialer.Dial(getWebSocketURL(validSearcherAddress, validCommitmentAddress), nil)
		assert.Equal(t, http.StatusForbidden, resp2.StatusCode)
		assert.Nil(t, conn2)
		assert.NotNil(t, resp2)
		assert.NotNil(t, err2)
	})

	// Test with a searcher that is already connected
	// TODO(@ckartik): Resolve this test as a searcher who tries to reconnect will fail
	t.Run("Already connected searcher is closes connection and reopens", func(t *testing.T) {
		validSearcherAddress := common.HexToAddress("0x812fC9524961d0566B3207fee1a567fef23E5E10")
		validCommitmentAddress := common.HexToAddress("0x812fC9524961d0566B3207fee1a567fef23E5E38")
		builderAddress := common.HexToAddress("0xbuilder")
		commitment := utils.GetCommitment(validCommitmentAddress, builderAddress)
		mockRollup := rollup.MockRollup{}
		mockRollup.On("GetStake", validSearcherAddress, commitment).Return(big.NewInt(100), nil)
		mockRollup.On("GetBuilderAddress").Return(builderAddress)
		mockRollup.On("GetMinimalStake", builderAddress).Return(big.NewInt(100))
		api.Rollup = &mockRollup

		conn, resp, err := dialer.Dial(getWebSocketURL(validSearcherAddress, validCommitmentAddress), nil)
		assert.NotNil(t, conn)
		assert.NotNil(t, resp)
		assert.Nil(t, err)
		conn.Close()

		_, resp2, err2 := dialer.Dial(getWebSocketURL(validSearcherAddress, validCommitmentAddress), nil)
		assert.Equal(t, http.StatusSwitchingProtocols, resp2.StatusCode)
		assert.NotNil(t, resp2)
		assert.Nil(t, err2)
	})
}
