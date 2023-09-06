package commons

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"golang.org/x/crypto/sha3"
)

var (
	nilAddress = common.Address{}
)

func BytesCompare(a, b []byte) bool {
	return bytes.Equal(a, b)
}

func PeersToBytes(p []peer.ID) []byte {
	out := bytes.NewBuffer(nil)
	for _, v := range p {
		out.Write([]byte(v.String()))
		out.Write([]byte(","))
	}
	return out.Bytes()
}

func GetNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// it takes a private key as input and returns the corresponding Ethereum address.
func GetAddressFromPrivateKey(privateKey *ecdsa.PrivateKey) common.Address {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("publicKeyECDSA assertion failed")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

// it takes compressed bytes representing a secp256k1 public key as input and returns
// the corresponding Ethereum address.
func Secp256k1CompressedBytesToEthAddress(compressedBytes []byte) common.Address {
	pb1, err := crypto.DecompressPubkey(compressedBytes)
	if err != nil {
		return nilAddress
	}

	pubkey := elliptic.Marshal(secp256k1.S256(), pb1.X, pb1.Y)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubkey[1:])
	address := hash.Sum(nil)

	return common.BytesToAddress(address[12:])
}

func ExtractIPFromMultiaddr(addr ma.Multiaddr) (net.IP, error) {
	// Use P_IP4 or P_IP6 protocol to extract the IP component
	ip4Bytes, err := addr.ValueForProtocol(ma.P_IP4)
	if err == nil {
		return net.ParseIP(ip4Bytes), nil
	}

	ip6Bytes, err := addr.ValueForProtocol(ma.P_IP6)
	if err == nil {
		return net.ParseIP(ip6Bytes), nil
	}

	// IP address not found
	return nil, fmt.Errorf("IP address not found in multiaddr")
}

func GetCallerName() string {
	pc, _, _, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	fullName := function.Name()
	parts := strings.Split(fullName, ".")
	return strings.ToLower(parts[len(parts)-1])
}
