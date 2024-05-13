package nanorpc

import (
	"context"

	"darvaza.org/core"
	"google.golang.org/protobuf/proto"
)

func (cs *clientConnectionState) unsafeFindRequestCallback(reqID int32) (int, bool) {
	if cs != nil {
		for i, x := range cs.cb {
			if x.RequestID == reqID {
				return i, true
			}
		}
	}

	return -1, false
}

func (c *Client) unsafeEnqueue(m *NanoRPCRequest, payload proto.Message, cb RequestCallback) error {
	switch {
	case c.cs == nil:
		// connecting
		return ErrNotConnected
	case c.cs.err != nil:
		// closing
		return c.cs.err
	}

	// encode
	msg, err := EncodeRequest(m, payload)
	if err != nil {
		return err
	}

	if cb != nil {
		// remember callback
		x := clientRequestQueue{
			RequestID:   m.RequestId,
			RequestType: m.RequestType,
			Callback:    cb,
		}
		c.cs.cb = append(c.cs.cb, x)
	}

	// enqueue
	c.cs.out <- [][]byte{msg}
	return nil
}

func (cs *clientConnectionState) unsafeDequeueAll() {
	fn := func(_ []clientRequestQueue, x clientRequestQueue) (clientRequestQueue, bool) {
		cb := x.Callback
		reqID := x.RequestID

		cs.Go(func(ctx context.Context) {
			_ = cb(ctx, reqID, nil)
		})

		return clientRequestQueue{}, false // discard
	}

	cs.cb = core.SliceReplaceFn(cs.cb, fn)
}

func (*clientConnectionState) handleResponse(_ context.Context, _ *NanoRPCResponse) error {
	return core.Wrap(core.ErrTODO, "handleResponse")
}
