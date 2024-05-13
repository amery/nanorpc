package nanorpc

import (
	"context"
	"net"
	"time"

	"darvaza.org/sidecar/pkg/reconnect"
	"darvaza.org/slog"
	"darvaza.org/slog/handlers/discard"
	"darvaza.org/x/config"
)

// NewWaitReconnect creates a function for [ClientOptions.WaitReconnect]
// that waits either for the cancellation of a [context.Context] or
// a specified duration.
func NewWaitReconnect(d time.Duration) func(context.Context) error {
	return reconnect.NewConstantWaiter(d)
}

// ClientOptions describes how the [Client] will operate
type ClientOptions struct {
	Context context.Context
	Logger  slog.Logger

	KeepAlive      time.Duration `default:"5s"`
	ConnectTimeout time.Duration `default:"2s"`
	ReadTimeout    time.Duration `default:"2s"`
	WriteTimeout   time.Duration `default:"2s"`

	ReconnectDelay time.Duration `default:"5s"`
	WaitReconnect  reconnect.Waiter

	AlwaysHashPaths bool
	HashCache       *HashCache

	// OnConnect is called when the connection is established.
	OnConnect func(context.Context, net.Conn) error
	// OnDisconnect is called after closing the connection and can be used to
	// prevent further connection retries.
	OnDisconnect func(context.Context, net.Conn) error
	// OnError is called after all errors and gives us the opportunity to
	// decide how the error should be treated by the reconnection logic.
	OnError func(context.Context, net.Conn, error) error
}

// SetDefaults fills gaps in [ClientOptions]
func (opts *ClientOptions) SetDefaults() error {
	if err := config.Set(opts); err != nil {
		return err
	}

	if opts.Context == nil {
		opts.Context = context.Background()
	}

	if opts.Logger == nil {
		opts.Logger = discard.New()
	}

	if opts.WaitReconnect == nil {
		opts.WaitReconnect = NewWaitReconnect(opts.ReconnectDelay)
	}

	if opts.HashCache == nil {
		// use global cache
		opts.HashCache = hashCache
	}

	return nil
}

// Export generates a [reconnect.Options]
func (opts *ClientOptions) Export() (*reconnect.Options, error) {
	if err := opts.SetDefaults(); err != nil {
		return nil, err
	}

	ro := &reconnect.Options{
		Context: opts.Context,
		Logger:  opts.Logger,

		KeepAlive:      opts.KeepAlive,
		ConnectTimeout: opts.ConnectTimeout,
		ReadTimeout:    opts.ReadTimeout,
		WriteTimeout:   opts.WriteTimeout,

		WaitReconnect: opts.WaitReconnect,
	}

	return ro, nil
}

func (opts *ClientOptions) getHashCache() *HashCache {
	if hc := opts.HashCache; hc != nil {
		// use given HashCache
		return hc
	}

	// use global cache
	return hashCache
}

func (opts *ClientOptions) newGetPathOneOf(hc *HashCache) func(string) isNanoRPCRequest_PathOneof {
	if opts.AlwaysHashPaths {
		// use path_hash
		if hc == nil {
			hc = opts.getHashCache()
		}

		return func(path string) isNanoRPCRequest_PathOneof {
			return &NanoRPCRequest_PathHash{
				PathHash: hc.Hash(path),
			}
		}
	}

	// use string
	return func(path string) isNanoRPCRequest_PathOneof {
		return &NanoRPCRequest_Path{
			Path: path,
		}
	}
}

// NewClient a new [Client] with default options
func NewClient(ctx context.Context, network, address string) (*Client, error) {
	opts := ClientOptions{
		Context: ctx,
	}

	return opts.New(network, address)
}

// New creates a new [Client] using given [ClientOptions].
func (opts *ClientOptions) New(network, address string) (*Client, error) {
	var c = new(Client)

	ro, err := opts.Export()
	if err != nil {
		return nil, err
	}

	rc, err := ro.New(network, address, c.preInit)
	if err != nil {
		return nil, err
	}

	if err := c.init(opts, rc); err != nil {
		return nil, err
	}

	return c, nil
}
