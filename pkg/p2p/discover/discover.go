package discover

import (
	"context"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/lthibault/log"
	"github.com/primev/builder-boost/pkg/p2p/commons"
	"github.com/primev/builder-boost/pkg/p2p/config"
)

// optional discovery
type Discovery struct {
	cfg    *config.Config
	h      host.Host
	ctx    context.Context
	log    log.Logger
	notify *discoveryNotifee
	idht   *dht.IpfsDHT

	//...
}

// Create discover method
func NewDiscovery(cfg *config.Config, h host.Host, ctx context.Context, logger log.Logger) *Discovery {
	notify := &discoveryNotifee{
		h:   h,
		ctx: ctx,
		log: logger,
	}

	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	idht, err := dht.New(ctx, h)
	if err != nil {
		logger.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p dht discover",
		}).Fatal(err)
	}
	if err = idht.Bootstrap(ctx); err != nil {
		logger.With(log.F{
			"caller":  commons.GetCallerName(),
			"date":    commons.GetNow(),
			"service": "p2p bootstrap discover",
		}).Fatal(err)
	}

	return &Discovery{
		cfg:    cfg,
		h:      h,
		ctx:    ctx,
		log:    logger,
		notify: notify,
		idht:   idht,
	}
}

// StartDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func (d *Discovery) StartMdnsDiscovery() error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(d.h, d.cfg.DiscoveryServiceTag(), d.notify)

	return s.Start()
}

func (d *Discovery) StartDhtRouting() {
	go d.discoverPeersWithRendezvous()
}

func (d *Discovery) ConnectBootstrap() {
	var wg sync.WaitGroup
	for _, peerAddr := range d.cfg.BootstrapPeers() {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := d.h.Connect(d.ctx, *peerinfo); err != nil {
				d.log.With(log.F{
					"caller":  commons.GetCallerName(),
					"date":    commons.GetNow(),
					"service": "p2p bootstrap",
				}).Error(err)

			} else {
				d.log.With(log.F{
					"caller":  commons.GetCallerName(),
					"date":    commons.GetNow(),
					"service": "p2p bootstrap",
					"peer":    *peerinfo,
				}).Info("connection established")
			}
		}()
	}

	wg.Wait()
	return
}

// works for a rendezvous peer discovery !!!
// TODO
func (d *Discovery) discoverPeersWithRendezvous() {

	routingDiscovery := drouting.NewRoutingDiscovery(d.idht)
	dutil.Advertise(d.ctx, routingDiscovery, d.cfg.Topic())

	// Look for others who have announced and attempt to connect to them
	d.log.With(log.F{
		"caller":  commons.GetCallerName(),
		"date":    commons.GetNow(),
		"service": "p2p rendezvous",
	}).Info("searching for peers...")

	for i := 0; i < 60; i++ {
		time.Sleep(time.Second)
		peerChan, err := routingDiscovery.FindPeers(d.ctx, d.cfg.Topic())
		if err != nil {
			d.log.With(log.F{
				"caller":  commons.GetCallerName(),
				"date":    commons.GetNow(),
				"service": "p2p rendezvous",
			}).Warn(err)
		}

		for peer := range peerChan {
			if peer.ID == d.h.ID() {
				continue // No self connection
			}
			err := d.h.Connect(d.ctx, peer)
			if err != nil {
				d.log.With(log.F{
					"caller":  commons.GetCallerName(),
					"date":    commons.GetNow(),
					"service": "p2p rendezvous",
					"peer":    peer.ID.Pretty(),
				}).Error(err)
			} else {
				d.notify.HandlePeerFound(peer)
				break
			}
		}
	}

	d.log.With(log.F{
		"caller":   commons.GetCallerName(),
		"date":     commons.GetNow(),
		"service":  "p2p rendezvous",
		"log time": commons.GetNow(),
	}).Info("peer discovery completed")

}
