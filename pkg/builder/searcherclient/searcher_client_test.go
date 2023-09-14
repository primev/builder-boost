package searcherclient_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/attestantio/go-builder-client/api/capella"
	"github.com/primev/builder-boost/pkg/builder/searcherclient"
)

type testSearcherC struct {
	id        string
	payloads  []searcherclient.SuperPayload
	heartbeat time.Time
	validity  int64
	closed    bool
}

func (t *testSearcherC) ID() string {
	return t.id
}

func (t *testSearcherC) Send(_ context.Context, payload searcherclient.SuperPayload) error {
	t.payloads = append(t.payloads, payload)
	return nil
}

func (t *testSearcherC) GetInfo() searcherclient.SearcherInfo {
	return searcherclient.SearcherInfo{
		ID:        t.id,
		Heartbeat: t.heartbeat,
		Validity:  t.validity,
	}
}

func (t *testSearcherC) Close() error {
	t.closed = true
	return nil
}

func TestSearcherClient(t *testing.T) {
	t.Parallel()

	t.Run("new and close", func(t *testing.T) {
		client := searcherclient.NewSearcherClient(newTestLogger(), false)

		err := client.Close()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("add and remove searcher", func(t *testing.T) {
		client := searcherclient.NewSearcherClient(newTestLogger(), false)
		t.Cleanup(func() {
			err := client.Close()
			if err != nil {
				t.Fatal(err)
			}
		})

		srch := &testSearcherC{
			id: "test",
		}

		client.AddSearcher(srch)
		client.RemoveSearcher("test")

		if !srch.closed {
			t.Fatal("expected searcher to be closed")
		}

		info := client.GetSeacherInfo()
		if len(info) != 0 {
			t.Fatalf("expected info to be empty, got %v", info)
		}
	})

	t.Run("get searcher info", func(t *testing.T) {
		client := searcherclient.NewSearcherClient(newTestLogger(), false)
		t.Cleanup(func() {
			err := client.Close()
			if err != nil {
				t.Fatal(err)
			}
		})

		srch := &testSearcherC{
			id:        "test",
			heartbeat: time.Now(),
			validity:  100,
		}

		client.AddSearcher(srch)

		info := client.GetSeacherInfo()
		if len(info) != 1 {
			t.Fatalf("expected info to have 1 entry, got %v", info)
		}

		if info[0].ID != srch.id {
			t.Fatalf("expected id to be %v, got %v", srch.id, info[0].ID)
		}

		if info[0].Heartbeat != srch.heartbeat {
			t.Fatalf("expected heartbeat to be %v, got %v", srch.heartbeat, info[0].Heartbeat)
		}

		if info[0].Validity != srch.validity {
			t.Fatalf("expected validity to be %v, got %v", srch.validity, info[0].Validity)
		}

		if !client.IsConnected(srch.id) {
			t.Fatalf("expected searcher to be connected")
		}
	})

	t.Run("send payload", func(t *testing.T) {
		client := searcherclient.NewSearcherClient(newTestLogger(), true)
		t.Cleanup(func() {
			err := client.Close()
			if err != nil {
				t.Fatal(err)
			}
		})

		srch := &testSearcherC{
			id: "0xeBc71Ae61372b940c940cF51510dEAB056E1f24C",
		}

		client.AddSearcher(srch)

		var payload capella.SubmitBlockRequest

		err := json.Unmarshal([]byte(jsonTxn), &payload)
		if err != nil {
			t.Fatal(err)
		}

		err = client.SubmitBlock(context.Background(), &payload)
		if err != nil {
			t.Fatal(err)
		}

		if len(srch.payloads) != 1 {
			t.Fatalf("expected payloads to have 1 entry, got %v", srch.payloads)
		}

		if len(srch.payloads[0].SearcherTxns) != 1 {
			t.Fatalf(
				"expected searcher txns to have 1 entry, got %v",
				srch.payloads[0].SearcherTxns,
			)
		}

		if srch.payloads[0].InternalMetadata.Transactions.Count != 1 {
			t.Fatalf(
				"expected internal metadata to have 1 entry, got %v",
				srch.payloads[0].InternalMetadata.Transactions.Count,
			)
		}

		if srch.payloads[0].InternalMetadata.Transactions.MaxPriorityFee != 25000000000 {
			t.Fatalf(
				"expected internal metadata to have max priority fee of 25000000000, got %v",
				srch.payloads[0].InternalMetadata.Transactions.MaxPriorityFee,
			)
		}

		if srch.payloads[0].InternalMetadata.Transactions.MinPriorityFee != 25000000000 {
			t.Fatalf(
				"expected internal metadata to have min priority fee of 25000000000, got %v",
				srch.payloads[0].InternalMetadata.Transactions.MaxPriorityFee,
			)
		}
	})
}

var jsonTxn = `{
				"message": {
					"slot": "2246842",
					"parent_hash": "0x04509e89bb0a974344d0855af9c6e41f26f3198c9edb1aac21e5c11e04745ef4",
					"block_hash": "0x2a4f3ec7cd9572a9b28aa7fb0a2c7799b90681b67f04cd64bd5e01023b81bc3b",
					"builder_pubkey": "0xaa1488eae4b06a1fff840a2b6db167afc520758dc2c8af0dfb57037954df3431b747e2f900fe8805f05d635e9a29717b",
					"proposer_pubkey": "0xa66d5b1cf24a38a598a45d16818d04e1c1331f8535591e7b9d3d13e390bfb466a0180098b4656131e087b72bf10be172",
					"proposer_fee_recipient": "0xf24a01ae29dec4629dfb4170647c4ed4efc392cd",
					"gas_limit": "30000000",
					"gas_used": "1884333",
					"value": "4429853700147000"
				},
				"execution_payload": {
					"parent_hash": "0x04509e89bb0a974344d0855af9c6e41f26f3198c9edb1aac21e5c11e04745ef4",
					"fee_recipient": "0x812fc9524961d0566b3207fee1a567fef23e5e38",
					"state_root": "0x1c70df2d36ef350d3cdacb781e0153d00c1f50d77c08ff4740b2185b741483c7",
					"receipts_root": "0xafd17a85f9de5eefa2f7e38b44a2fc560ede4dbf81433d3dbbf7da6560e49bcc",
					"logs_bloom": "0x0000020000100000000002040802000001012408000024000000000000000008000a00000004040000800000148011000008100000001000080000108020000000002008c00020000080000e600010004080100440802000000000000001000000698000060018020084200205100940000002900900000100000810000000004100040000000000000040084004800001000000003100000004000000000000020800000000000100014100400802002040000000000400000008000000400400000022200000004400000000004000200000480200000000000100000020001034000000002000002000002000000080000000080800001204000000106000",
					"prev_randao": "0xa868ab4cfafc9107183622c8d955b293558cfbaa661fd1ae41403db3f3b35618",
					"block_number": "3379258",
					"gas_limit": "30000000",
					"gas_used": "1884333",
					"timestamp": "1682695704",
					"extra_data": "0xd883010b05846765746888676f312e32302e33856c696e7578",
					"base_fee_per_gas": "7",
					"block_hash": "0x2a4f3ec7cd9572a9b28aa7fb0a2c7799b90681b67f04cd64bd5e01023b81bc3b",
					"transactions": [
						"0x02f901b883aa36a782538c8505d21dba00850ba43b74008316e360944a32aa121683ba185e2a55a52a9881135ab91e9a80b90144979986110000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000400000000000000000000000092c55b159f45648957f32c8a017ac7d62b16e1f744cd496a1dc36bcd4a73510cb07ac0a0c9989c3dbdfa6f56b5d54307286956a900000000000000000000000092c55b159f45648957f32c8a017ac7d62b16e1f7ac8b9b632aec60ff00f2e8877aeb2f3fc22842ea5f11ccffdac7db0d60d6655500000000000000000000000092c55b159f45648957f32c8a017ac7d62b16e1f76109664ef022c6a16e05d84e1f07416863fa669258886f750beb9d930e3be22f00000000000000000000000092c55b159f45648957f32c8a017ac7d62b16e1f78ddb09b000ddbb3cc50d7b7ebd90b69eb87a340b4138f9ffcc2e5e1831bfa955c080a03024de7fe25e0c43fd0bbd649a5616de42734ac627fd1e67c9424f66a99aa809a004ecfed7c69b0c629e1e7b54c0892aa44589a7ee6f19ae58949ab48b11a1b9c7"
					],
					"withdrawals": []
				},
				"signature": "0x890edd3019111f248105b150d41e31b2a526c01ba6d8695943394796e5825c34b693efdcb05e89b0ec544d08bae6463c098e5e3855b03b1a707398cccb2c7168c86233816f6a4a026fc0d0ee660121f0951d8a76dc3c9b3dcdd7a832e0f5c44d"
			}`
