package preconf

// construct test for bid
import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestBid(t *testing.T) {
	key, _ := crypto.GenerateKey()
	bid, err := ConstructSignedBid(big.NewInt(10), "0xkartik", big.NewInt(2), key)
	if err != nil {
		t.Fatal(err)
	}
	address, err := bid.VerifySearcherSignature()
	if err != nil {
		t.Fatal(err)
	}

	if address.Big().Cmp(crypto.PubkeyToAddress(key.PublicKey).Big()) != 0 {
		t.Fatal("Address not same as signer")
	}
}

func TestCommitment(t *testing.T) {
	key, _ := crypto.GenerateKey()
	bid, err := ConstructSignedBid(big.NewInt(10), "0xkartik", big.NewInt(2), key)
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

}
