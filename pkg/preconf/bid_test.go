package preconf

// construct test for bid
import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestBid(t *testing.T) {
	// Connect to Ganache
	key, _ := crypto.GenerateKey()

	bid, err := ConstructSignedBid(big.NewInt(10), "0xkartik", big.NewInt(2), key)
	if err != nil {
		t.Fatal(err)
	}
	address, err := bid.VerifySearcherSignature()
	t.Log(address)
	b, _ := json.Marshal(bid)
	var bid2 PreConfBid
	json.Unmarshal(b, &bid2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hexutil.Bytes(bid2.BidHash))
	sig := hexutil.Bytes(bid2.Signature)
	t.Log(sig)

	if address.Big().Cmp(crypto.PubkeyToAddress(key.PublicKey).Big()) != 0 {
		t.Fatal("Address not same as signer")
	}
}

func TestCommitment(t *testing.T) {
	client, err := ethclient.Dial("http://54.200.76.18:8545")
	if err != nil {
		t.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	key, _ := crypto.GenerateKey()
	bid, err := ConstructSignedBid(big.NewInt(10), "0xkadrtik", big.NewInt(2), key)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.Marshal(bid)
	var bid2 PreConfBid
	json.Unmarshal(b, &bid2)
	commit, err := bid2.ConstructCommitment(key)

	if err != nil {
		t.Fatal(err)
	}
	commit.VerifyBuilderSignature()
	privateKey, _ := crypto.HexToECDSA("a9b36f394d2133174158d2eed84ffd4da979a73fd26eaa7a516fe4927ec29bcc")

	txn, err := commit.StoreCommitmentToDA(privateKey, "0xac27A2cbdBA8768D49e359ebA326fC1F27832ED4", client)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txn.Hash().Hex())
}
