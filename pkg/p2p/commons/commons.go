package commons

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
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

func GetAddressFromPrivateKey(privateKey *ecdsa.PrivateKey) common.Address {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("publicKeyECDSA assertion failed")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
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
