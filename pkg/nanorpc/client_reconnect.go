package nanorpc

import (
	"context"
	"net"

	"darvaza.org/sidecar/pkg/reconnect"
)

// preInit adjusts the [reconnect.Options] to use
// [Client]'s callbacks.
func (c *Client) preInit(opts *reconnect.Options) error {
	opts.OnSession = c.reconnectSession
	opts.OnDisconnect = c.reconnectDisconnect
	opts.OnError = c.reconnectError
	return nil
}

func (c *Client) reconnectSession(ctx context.Context) error {
	c.setup(ctx)
	return c.cs.Run(ctx)
}

func (c *Client) reconnectDisconnect(ctx context.Context, conn net.Conn) error {
	c.rc.SayRemote(conn, "%s: %s", "nanorpc", "disconnected")

	// notify any pending callback
	c.mu.Lock()
	if cs := c.cs; cs != nil {
		cs.unsafeDequeueAll()
	}
	c.mu.Unlock()

	// notify the user
	if fn := c.onDisconnect; fn != nil {
		return fn(ctx, conn)
	}

	return nil
}

func (c *Client) reconnectError(_ context.Context, conn net.Conn, err error) error {
	if reconnect.IsNoiseError(err) {
		return nil
	}

	c.rc.SayRemoteError(conn, err, "%s: %s (%T)", "nanorpc", "error", err)
	return err
}

// Connect initiates the self-reconnecting NanoRPC client.
func (c *Client) Connect() error { return c.rc.Connect() }

// Wait blocks until the worker launched by [Connect]
// has finished.
func (c *Client) Wait() error { return c.rc.Wait() }

// Close initiates a shutdown.
func (c *Client) Close() error { return c.rc.Close() }

// Shutdown initiates a shutdown and wait until it's all done,
// or the given context reaches the deadline.
func (c *Client) Shutdown(ctx context.Context) error { return c.rc.Shutdown(ctx) }
