package preconf

// construct test for bid
import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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
	privateKey, _ := crypto.HexToECDSA("2c440a842f6afc128360dc3e7cafefd66eef4929e7df50bdf0372b93f524ea63")

	txn, err := commit.StoreCommitmentToDA(privateKey, "0x1D2D5eC8afeF61AfC9B25bEd51681A3B80989473", "http://localhost:8545")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txn.Hash().Hex())
}
