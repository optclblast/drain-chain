package tcp

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/optclblast/drain-chain/internal/kernel/logger"
)

type TCPInterface struct {
	log logger.Logger
	l   net.Listener

	closed_m     sync.Mutex
	closed       bool
	connLifetime time.Duration

	conns_m sync.Mutex
	conns   map[string]net.Conn
}

func NewTCPInterface(
	log logger.Logger,
	bindAddr string,
	connLifetime time.Duration,
) (*TCPInterface, func(), error) {
	l, err := net.Listen("tcp", bindAddr)
	if err != nil {
		return nil, func() {}, fmt.Errorf("system.tcp: error listent to %s. %w", bindAddr, err)
	}

	t := &TCPInterface{
		log:          log,
		l:            l,
		connLifetime: connLifetime,
		conns:        make(map[string]net.Conn),
	}

	return t, func() {
		l.Close()
	}, nil
}

func (t *TCPInterface) Start() {
	go t.acceptConnections()
}

func (t *TCPInterface) Stop(_ context.Context) error {
	t.closed_m.Lock()
	defer t.closed_m.Unlock()

	if t.closed {
		return nil
	}

	t.closed = true

	return nil
}

func (t *TCPInterface) Reload(ctx context.Context) {
	t.Stop(ctx)
	t.Start()
}

func (t *TCPInterface) acceptConnections() {
	for !t.closed {
		c, err := t.l.Accept()
		if err != nil {
			t.log.Error(
				"system.tcp.acceptLoop: error accept connection",
				logger.Err(err),
			)

			continue
		}

		t.conns_m.Lock()
		t.conns[c.RemoteAddr().String()] = c
		t.conns_m.Unlock()
	}
}
