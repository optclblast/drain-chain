package p2p

import (
	"context"
	"net"
	"sync"

	"github.com/optclblast/drain-chain/internal/kernel/logger"
	"github.com/optclblast/drain-chain/internal/kernel/otp/application"
	"github.com/optclblast/drain-chain/internal/kernel/system/network/tcp"
)

type Peer struct {
	Address string
}

// P2P is an application, that is responsible for p2p communications of the node
type P2P interface {
	application.Application

	Connect(ctx context.Context, peer Peer) (net.Conn, error)
}

type p2p struct {
	closed   bool
	closed_m sync.Mutex

	address string
	tcp     *tcp.TCPInterface
	log     *logger.Logger
}

func NewP2P(
	address string,
	tcp *tcp.TCPInterface,
	log *logger.Logger,
) P2P {
	return &p2p{
		address: address,
		tcp:     tcp,
		log:     log.WithGroup("system.network.p2p"),
	}
}

func (t *p2p) Connect(ctx context.Context, peer Peer) (net.Conn, error) {
	return nil, nil
}

func (t *p2p) Start() {
	// todo start app here
}

func (t *p2p) Stop(_ context.Context) error {
	t.closed_m.Lock()
	defer t.closed_m.Unlock()

	if t.closed {
		return nil
	}

	t.closed = true

	return nil
}

func (t *p2p) Reload(ctx context.Context) {
	t.Stop(ctx)
	t.Start()
}
