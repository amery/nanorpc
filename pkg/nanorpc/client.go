package nanorpc

import (
	"context"
	"net"
	"sync"

	"darvaza.org/sidecar/pkg/reconnect"
)

// Client is a reconnecting NanoRPC client.
type Client struct {
	mu sync.Mutex
	rc *reconnect.Client
	hc *HashCache
	cs *clientConnectionState

	reqCounter   int32
	getPathOneOf func(string) isNanoRPCRequest_PathOneof

	onConnect    func(context.Context, net.Conn) error
	onDisconnect func(context.Context, net.Conn) error
	onError      func(context.Context, net.Conn, error) error
}

func (c *Client) init(opts *ClientOptions, rc *reconnect.Client) error {
	c.rc = rc
	c.hc = opts.getHashCache()
	c.getPathOneOf = opts.newGetPathOneOf(c.hc)

	c.onConnect = opts.OnConnect
	c.onDisconnect = opts.OnDisconnect
	c.onError = opts.OnError

	return nil
}

// RequestCallback handles a response to a request
type RequestCallback func(context.Context, int32, *NanoRPCResponse) error
