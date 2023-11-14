package reconnect

import (
	"context"
	"net"
	"time"

	"darvaza.org/core"
)

// A ClientOption alters how the [Client] is initially configured
type ClientOption func(*Client) error

// WithContext specifies the parent context the Client's worker
// will use for cancellations.
func WithContext(ctx context.Context) ClientOption {
	return func(conn *Client) error {
		switch {
		case ctx == nil:
			return core.Wrap(core.ErrInvalid, "no context specified")
		case conn.ctx != nil:
			return core.Wrap(core.ErrExists, "client already has a parent context")
		default:
			conn.ctx = ctx
			return nil
		}
	}
}

// WithDialer provides an optional dialer to establish connections
func WithDialer(dialer *net.Dialer) ClientOption {
	return func(conn *Client) error {
		// just in case someone decides to call this after
		// the [Client] is created.
		conn.mu.Lock()
		defer conn.mu.Unlock()

		switch {
		case dialer == nil:
			return core.Wrap(core.ErrInvalid, "no dialer specfied")
		default:
			conn.dialer = dialer
			return nil
		}
	}
}

// WithStartWait specifies how long will we give the connection worker
// for an early error.
func WithStartWait(d time.Duration) ClientOption {
	return func(conn *Client) error {
		switch {
		case d < time.Millisecond:
			return core.Wrap(core.ErrInvalid, "invalid health wait (%s)", d)
		default:
			conn.waitStart = d
			return nil
		}
	}
}

func WithNotifier(notifiers ...NotifierFunc) ClientOption {
	return func(conn *Client) error {
		// just in case someone decides to call this after
		// the [Client] is created.
		conn.mu.Lock()
		defer conn.mu.Unlock()

		conn.notifiers = append(conn.notifiers, notifiers...)
		return nil
	}
}
