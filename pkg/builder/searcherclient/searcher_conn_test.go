package searcherclient_test

import (
	"context"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/primev/builder-boost/pkg/builder/searcherclient"
	"golang.org/x/exp/slog"
)

var upgrader = websocket.Upgrader{}

type wsMsg struct {
	msgType int
	msg     []byte
}

type testSearcher struct {
	inbox chan wsMsg
}

func (t *testSearcher) handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()

	c.SetCloseHandler(func(code int, text string) error {
		t.inbox <- wsMsg{msgType: websocket.CloseMessage, msg: []byte(text)}
		return nil
	})

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		t.inbox <- wsMsg{msgType: mt, msg: message}
	}
}

func newTestServer(t *testing.T) (*testSearcher, *websocket.Conn) {
	s := &testSearcher{
		inbox: make(chan wsMsg),
	}

	srv := httptest.NewServer(http.HandlerFunc(s.handler))

	t.Cleanup(func() {
		srv.Close()
	})

	// Convert http://127.0.0.1 to ws://127.0.0.1
	u := "ws" + strings.TrimPrefix(srv.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Cleanup(func() {
		_ = ws.Close()
	})

	return s, ws
}

type testRollup struct {
	blockno int64
}

func (r *testRollup) GetBlockNumber() (*big.Int, error) {
	return big.NewInt(r.blockno), nil
}

func newTestLogger() *slog.Logger {
	testLogger := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(testLogger)
}

func TestSearcherConn(t *testing.T) {
	t.Parallel()

	t.Run("new and close", func(t *testing.T) {
		s, ws := newTestServer(t)
		rollUp := &testRollup{blockno: 0}

		var disconnectReason string

		// Create a new searcher
		srch := searcherclient.NewSearcher(
			"test",
			ws,
			rollUp,
			newTestLogger(),
			100,
			func(reason string) {
				disconnectReason = reason
			},
		)

		// Close the searcher
		err := srch.Close()
		if err != nil {
			t.Fatal(err)
		}

		msg := <-s.inbox
		if msg.msgType != websocket.CloseMessage {
			t.Fatalf("expected close message, got %v", msg.msgType)
		}
		if disconnectReason != "server closing connection" {
			t.Fatalf(
				"expected disconnect reason to be 'server closing connection', got %v",
				disconnectReason,
			)
		}
	})

	t.Run("get searcher info after init", func(t *testing.T) {
		_, ws := newTestServer(t)
		rollUp := &testRollup{blockno: 0}

		// Create a new searcher
		srch := searcherclient.NewSearcher(
			"test",
			ws,
			rollUp,
			newTestLogger(),
			100,
			func(reason string) {},
		)

		t.Cleanup(func() {
			err := srch.Close()
			if err != nil {
				t.Errorf("failed to close searcher: %v", err)
			}
		})

		info := srch.GetInfo()
		if srch.ID() != info.ID {
			t.Fatalf("expected id to be %v, got %v", srch.ID(), info.ID)
		}
		if info.Heartbeat.Unix() == 0 {
			t.Fatalf("expected heartbeat to be non-zero")
		}
		if info.Validity != 100 {
			t.Fatalf("expected validity to be 100, got %v", info.Validity)
		}
	})

	t.Run("send payload", func(t *testing.T) {
		s, ws := newTestServer(t)
		rollUp := &testRollup{blockno: 0}

		// Create a new searcher
		srch := searcherclient.NewSearcher(
			"test",
			ws,
			rollUp,
			newTestLogger(),
			100,
			func(reason string) {},
		)

		t.Cleanup(func() {
			err := srch.Close()
			if err != nil {
				t.Errorf("failed to close searcher: %v", err)
			}
		})

		info := srch.GetInfo()
		time.Sleep(1 * time.Second)

		err := srch.Send(context.Background(), searcherclient.SuperPayload{
			InternalMetadata: searcherclient.Metadata{
				Builder:   "test",
				Number:    0,
				BlockHash: "test",
			},
			SearcherTxns: map[string][]string{
				"test": {"test"},
			},
		})
		if err != nil {
			t.Fatalf("failed to send payload: %v", err)
		}

		msg := <-s.inbox
		if msg.msgType != websocket.TextMessage {
			t.Fatalf("expected text message, got %v", msg.msgType)
		}

		var payload searcherclient.SuperPayload
		err = json.Unmarshal(msg.msg, &payload)
		if err != nil {
			t.Fatalf("failed to unmarshal payload: %v", err)
		}

		if payload.InternalMetadata.Builder != "test" {
			t.Fatalf("expected builder to be 'test', got %v", payload.InternalMetadata.Builder)
		}

		if payload.InternalMetadata.Number != 0 {
			t.Fatalf("expected number to be 0, got %v", payload.InternalMetadata.Number)
		}

		if payload.InternalMetadata.BlockHash != "test" {
			t.Fatalf("expected block hash to be 'test', got %v", payload.InternalMetadata.BlockHash)
		}

		if len(payload.SearcherTxns) != 1 {
			t.Fatalf("expected 1 searcher txns, got %v", len(payload.SearcherTxns))
		}

		if len(payload.SearcherTxns["test"]) != 1 {
			t.Fatalf(
				"expected 1 searcher txns for 'test', got %v",
				len(payload.SearcherTxns["test"]),
			)
		}

		if payload.SearcherTxns["test"][0] != "test" {
			t.Fatalf("expected searcher txn to be 'test', got %v", payload.SearcherTxns["test"][0])
		}

		newInfo := srch.GetInfo()
		if !newInfo.Heartbeat.After(info.Heartbeat) {
			t.Fatalf("expected id to be %v, got %v", info, newInfo)
		}
	})

	t.Run("subscription ended", func(t *testing.T) {
		s, ws := newTestServer(t)
		rollUp := &testRollup{blockno: 101}
		disconnectReason := ""

		// Create a new searcher
		srch := searcherclient.NewSearcher(
			"test",
			ws,
			rollUp,
			newTestLogger(),
			100,
			func(reason string) { disconnectReason = reason },
		)

		t.Cleanup(func() {
			err := srch.Close()
			if err != nil {
				t.Errorf("failed to close searcher: %v", err)
			}
		})

		err := srch.Send(context.Background(), searcherclient.SuperPayload{
			InternalMetadata: searcherclient.Metadata{
				Builder:   "test",
				Number:    0,
				BlockHash: "test",
			},
			SearcherTxns: map[string][]string{
				"test": {"test"},
			},
		})
		if err != nil {
			t.Fatalf("failed to send payload: %v", err)
		}

		msg := <-s.inbox
		if msg.msgType != websocket.CloseMessage {
			t.Fatalf("expected Close message, got %v", msg.msgType)
		}

		if disconnectReason != "subscription ended" {
			t.Fatalf(
				"expected disconnect reason to be 'subscription ended', got %v",
				disconnectReason,
			)
		}
	})

	t.Run("client closed connection", func(t *testing.T) {
		_, ws := newTestServer(t)
		rollUp := &testRollup{blockno: 0}
		disconnectReason := ""
		closed := make(chan struct{})

		// Create a new searcher
		srch := searcherclient.NewSearcher(
			"test",
			ws,
			rollUp,
			newTestLogger(),
			100,
			func(reason string) { disconnectReason = reason; close(closed) },
		)

		t.Cleanup(func() {
			err := srch.Close()
			if err != nil {
				t.Errorf("failed to close searcher: %v", err)
			}
		})

		err := ws.Close()
		if err != nil {
			t.Fatalf("failed to close websocket: %v", err)
		}

		err = srch.Send(context.Background(), searcherclient.SuperPayload{
			InternalMetadata: searcherclient.Metadata{
				Builder:   "test",
				Number:    0,
				BlockHash: "test",
			},
			SearcherTxns: map[string][]string{
				"test": {"test"},
			},
		})
		if err != nil {
			t.Fatalf("failed to send payload: %v", err)
		}

		<-closed

		if disconnectReason != "failed to send payload" {
			t.Fatalf(
				"expected disconnect reason to be 'failed to send', got %v",
				disconnectReason,
			)
		}
	})
}
