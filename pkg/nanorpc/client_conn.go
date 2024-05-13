package nanorpc

import (
	"bufio"
	"context"
	"sync"

	"darvaza.org/core"
	"darvaza.org/sidecar/pkg/reconnect"
)

type clientConnectionState struct {
	c      *Client
	rc     *reconnect.Client
	ctx    context.Context
	cancel context.CancelCauseFunc
	err    error
	wg     sync.WaitGroup

	cb  []clientRequestQueue
	in  chan *NanoRPCResponse
	out chan [][]byte
}

type clientRequestQueue struct {
	RequestID   int32
	RequestType NanoRPCRequest_Type
	Callback    RequestCallback
}

func (c *Client) setup(parent context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ctx, cancel := context.WithCancelCause(parent)

	// init
	cs := &clientConnectionState{
		c:      c,
		rc:     c.rc,
		ctx:    ctx,
		cancel: cancel,

		in:  make(chan *NanoRPCResponse),
		out: make(chan [][]byte),
	}

	// store
	c.cs = cs

	// spawn reader
	cs.Go(func(ctx context.Context) {
		if err := cs.connReader(ctx); err != nil {
			cs.terminate(err)
		}
	})

	// spawn writer
	cs.Go(func(ctx context.Context) {
		if err := cs.connWriter(ctx); err != nil {
			cs.terminate(err)
		}
	})
}

func (cs *clientConnectionState) terminate(cause error) {
	cs.c.mu.Lock()
	defer cs.c.mu.Unlock()

	if cs.err == nil {
		// once
		if cause == nil {
			cause = context.Canceled
		}

		cs.err = cause
		cs.cancel(cause)
		cs.rc.Close()

		close(cs.in)
	}
}

func (cs *clientConnectionState) Go(fn func(context.Context)) {
	cs.wg.Add(1)
	go func() {
		defer cs.wg.Done()
		fn(cs.ctx)
	}()
}

func (cs *clientConnectionState) Err() error {
	cs.c.mu.Lock()
	defer cs.c.mu.Unlock()

	return cs.err
}

func (cs *clientConnectionState) Wait() error {
	cs.wg.Wait()
	return cs.Err()
}

func (cs *clientConnectionState) Run(ctx context.Context) error {
	cs.runLoop(ctx)
	return cs.Err()
}

func (cs *clientConnectionState) runLoop(ctx context.Context) {
	defer cs.wg.Wait()

	defer func() {
		if err := core.AsRecovered(recover()); err != nil {
			cs.terminate(err)
		}
	}()

	for {
		if err := cs.runOnce(ctx); err != nil {
			cs.terminate(err)
			break
		}
	}
}

func (cs *clientConnectionState) runOnce(ctx context.Context) error {
	select {
	case <-ctx.Done():
		// externally cancelled
		return ctx.Err()
	case <-cs.ctx.Done():
		// terminated by worker
		return cs.ctx.Err()
	case m := <-cs.in:
		if m != nil {
			return cs.handleResponse(ctx, m)
		}
	}

	return nil
}

func (cs *clientConnectionState) connReader(_ context.Context) error {
	s := bufio.NewScanner(cs.rc)
	s.Split(Split)

	for s.Scan() {
		msg, _, err := DecodeResponse(s.Bytes())
		if err != nil {
			return err
		}

		cs.in <- msg
	}

	return s.Err()
}

func (cs *clientConnectionState) connWriter(_ context.Context) error {
	for _, p := range <-cs.out {
		if _, err := cs.rc.Write(p); err != nil {
			return err
		}

		if err := cs.rc.Flush(); err != nil {
			return err
		}
	}

	return nil
}
