package config

import (
	"github.com/multiformats/go-multiaddr"
)

var DefaultBootstrapPeers []multiaddr.Multiaddr

func init() {
	for _, s := range []string{
		"/ip4/192.168.1.240/tcp/41859/p2p/12D3KooWRaKUthj9X6f3eLjyvjJp7DDAPo46NVi82UQcafMpMQKW",
		"/ip4/192.168.1.240/tcp/46771/p2p/12D3KooWBaAGA9xz5qnPAGxWjtZHgUNHB1jndM5Jp3mJQNTym1iE",
		//"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		//"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ", // mars.i.ipfs.io
	} {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			panic(err)
		}
		DefaultBootstrapPeers = append(DefaultBootstrapPeers, ma)
	}
}
