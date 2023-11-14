package reconnect

import (
	"context"
	"net"
	"sync"
	"syscall"
	"time"

	"darvaza.org/core"
)

// A Client asynchronously connects to a server,
// and reconnects if needed.
type Client struct {
	network   string
	address   string
	reconnect bool
	waitStart time.Duration

	dialer *net.Dialer
	conn   net.Conn
	ctx    context.Context
	cancel context.CancelFunc

	closed bool
	err    error
	mu     sync.Mutex
	wg     sync.WaitGroup

	notifiers []NotifierFunc
}

type Event uint
type NotifierFunc func(*Client, Event)

func (conn *Client) setDefaults() error {
	if conn.dialer == nil {
		conn.dialer = new(net.Dialer)
	}

	if conn.ctx == nil {
		conn.ctx = context.Background()
	}

	return nil
}

func (conn *Client) Err() error {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	return conn.err
}

func (conn *Client) storeError(err error) {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	conn.err = err
}

func (conn *Client) unsafeIsConnected() bool {
	return conn.conn != nil
}

// Connect initiates a self-retrying connection
func (conn *Client) Connect() error {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	if conn.unsafeIsConnected() {
		return syscall.EBUSY
	}

	conn.reset(true)
	return conn.unsafeStart(nil)
}

// Dial attempts a single connection
func (conn *Client) Dial() error {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	if conn.unsafeIsConnected() {
		return syscall.EBUSY
	}

	c, err := conn.dialer.DialContext(conn.ctx, conn.network, conn.address)
	if err != nil {
		return err
	}

	conn.reset(false)
	return conn.unsafeStart(c)
}

func (conn *Client) reset(reconnect bool) {
	conn.err = nil
	conn.closed = false
	conn.reconnect = reconnect
}

// Close closes the corruent connection and disable
// automatic reconnections.
// Use [Client.Connect] or [Client.Dial] to resume.
func (conn *Client) Close() error {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	conn.closed = true
	conn.reconnect = false

	if conn.unsafeIsConnected() {
		// disconnect
		err := conn.unsafeDisconnect()

		// and wait for the workers to finish
		conn.wg.Wait()
		return err
	}

	return nil
}

func (*Client) unsafeDisconnect() error {
	return core.ErrNotImplemented
}

func (conn *Client) unsafeStart(c net.Conn) error {
	var err error

	ctx, cancel := context.WithCancel(conn.ctx)
	conn.cancel = cancel

	conn.wg.Add(1)
	go func() {
		defer conn.wg.Done()

		defer func() {
			if rcv := recover(); rcv != nil {
				// panic
				err = core.NewPanicError(0, rcv)
				// if unsafeStart is waiting, err is already set
				// and will be returned to the caller.
				// otherwise, we wait for the Client to be unlocked
				conn.storeError(err)
			}
		}()

		err = conn.run(ctx, c)
		if err != nil {
			// if unsafeStart is waiting, err is already set
			// and will be returned to the caller.
			// otherwise, we wait for the Client to be unlocked
			conn.storeError(err)
		}
	}()

	if d := conn.waitStart; d > 0 {
		// wait for early failures
		<-time.After(d)
	}

	return err
}

func (*Client) run(_ context.Context, _ net.Conn) error {
	return core.ErrNotImplemented
}

// NewClient creates a new network [Client] to the specified server
// and optionally special configuration details.
func NewClient(network, address string, options ...ClientOption) (*Client, error) {
	// TODO: validate network and address
	conn := &Client{
		network: network,
		address: address,
	}

	for _, opt := range options {
		if err := opt(conn); err != nil {
			return nil, err
		}
	}

	if err := conn.setDefaults(); err != nil {
		return nil, err
	}

	return conn, nil
}
