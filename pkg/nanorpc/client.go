package nanorpc

import (
	"context"
	"net"
	"sync"
	"syscall"
	"time"

	"darvaza.org/core"
	"darvaza.org/slog"
	"darvaza.org/x/config"
)

// NewWaitReconnect creates a function for [ClientOptions.WaitReconnect]
// that waits either for the cancellation of a [context.Context] or
// a specified duration.
func NewWaitReconnect(d time.Duration) func(context.Context) error {
	return func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(d):
			return nil
		}
	}
}

// ClientOptions specifies how the NanoRPC client will operate.
type ClientOptions struct {
	Context context.Context
	Logger  slog.Logger

	KeepAlive      time.Duration `default:"5s"`
	ConnectTimeout time.Duration `default:"2s"`
	ReadTimeout    time.Duration `default:"2s"`
	WriteTimeout   time.Duration `default:"2s"`

	ReconnectDelay time.Duration
	WaitReconnect  func(ctx context.Context) error
}

// SetDefaults fills gaps in [ClientOptions].
func (opts *ClientOptions) SetDefaults() error {
	if err := config.Set(opts); err != nil {
		return err
	}

	if opts.Context == nil {
		opts.Context = context.Background()
	}

	if d := opts.ReconnectDelay; d > 0 && opts.WaitReconnect == nil {
		opts.WaitReconnect = NewWaitReconnect(d)
	}

	return nil
}

// New creates a new [Client] using given [ClientOptions].
func (opts *ClientOptions) New(network, address string) (*Client, error) {
	if err := opts.SetDefaults(); err != nil {
		return nil, err
	}

	// TODO: validate network and address

	c := &Client{
		options: *opts,
		dialer: net.Dialer{
			Timeout:   opts.ConnectTimeout,
			KeepAlive: opts.KeepAlive,
		},

		network: network,
		address: address,
	}

	c.ctx, c.cancel = context.WithCancelCause(opts.Context)

	return c, nil
}

// Client is a reconnecting NanoRPC client.
type Client struct {
	mu sync.Mutex

	options ClientOptions
	ctx     context.Context
	cancel  context.CancelCauseFunc
	err     error

	dialer  net.Dialer
	network string
	address string

	running bool
	conn    net.Conn
}

func (c *Client) dial() (net.Conn, error) {
	conn, err := c.dialer.DialContext(c.ctx, c.network, c.address)
	if err != nil {
		c.sayError(nil, err)
	}
	if conn != nil {
		c.say(conn, "connected")
	}
	return conn, err
}

func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return syscall.EBUSY
	}

	conn, err := c.dial()
	if err != nil && IsTimeout(err) {
		return err
	}

	c.running = true
	go c.run(conn)

	return nil
}

func (c *Client) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.err
}

func (c *Client) run(conn net.Conn) {
	defer func() {
		// panic?
		if err := core.AsRecovered(recover()); err != nil {
			c.setDone(err)
		}
	}()

	if conn != nil {
		// connected
		if err := c.setup(conn); err != nil {
			c.setDone(err)
			return
		}
	}

	for {
		var err error

		if c.conn == nil {
			err = c.reconnect()
		} else {
			err = c.runOnce()
		}

		if err != nil {
			c.setDone(err)
			return
		}
	}
}

func (c *Client) reconnect() error {
	wait := c.options.WaitReconnect

	if wait != nil {
		if err := wait(c.ctx); err != nil {
			// abort
			return err
		}
	}

	conn, err := c.dial()
	switch {
	case err == nil:
		return c.setup(conn)
	case IsTimeout(err) || IsTemporary(err):
		// try again later
		return nil
	default:
		// abort
		return err
	}
}

func (c *Client) setup(net.Conn) error
func (c *Client) runOnce() error

func (c *Client) setDone(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if conn := c.conn; conn != nil {
		defer conn.Close()

		c.conn = nil
	}

	c.err = err
	c.running = false
}

// NewClient a new [Client] with default options
func NewClient(ctx context.Context, network, address string) (*Client, error) {
	opts := ClientOptions{
		Context: ctx,
	}

	return opts.New(network, address)
}
