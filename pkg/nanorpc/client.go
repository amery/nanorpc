package nanorpc

import (
	"context"
	"net"
	"time"

	"github.com/amery/nanorpc/pkg/reconnect"
)

type Client struct {
	c *reconnect.Client
}

func WithContext(ctx context.Context) reconnect.ClientOption {
	return reconnect.WithContext(ctx)
}

func WithDialer(dialer *net.Dialer) reconnect.ClientOption {
	return reconnect.WithDialer(dialer)
}

func (c *Client) Connect() error {
	return c.c.Connect()
}

func (c *Client) Dial() error {
	return c.c.Dial()
}

func (c *Client) Close() error {
	return c.c.Close()
}

func (*Client) SetAutoReconnect(_ time.Duration) {
	// TODO: implement c.c.SetAutoReconnect
}

func (*Client) reconnectNotify(_ *reconnect.Client, _ reconnect.Event) {
	// TODO: do something useful
}

func NewClient(network, address string, options ...reconnect.ClientOption) (*Client, error) {
	var opts []reconnect.ClientOption

	c := new(Client)

	// we are first to be notified, and have the final word on startWait
	opts = append(opts, reconnect.WithNotifier(c.reconnectNotify))
	opts = append(opts, options...)
	opts = append(opts, reconnect.WithStartWait(time.Second))

	conn, err := reconnect.NewClient(network, address, opts...)
	if err != nil {
		return nil, err
	}
	c.c = conn

	return c, nil
}
