// Package main implements `nanorpc-cli`
package main

import (
	"context"
	"log"
	"time"

	"github.com/amery/nanorpc/pkg/nanorpc"

	"darvaza.org/sidecar/pkg/reconnect"
	"darvaza.org/slog"
)

const (
	// DefaultReconnectDelay specifies how long we wait
	// between attempts by default.
	DefaultReconnectDelay = 10 * time.Second
)

func main() {
	var ctx, cancel = context.WithCancelCause(context.Background())
	defer cancel(nil)

	var conf = &nanorpc.ClientOptions{
		Context:        ctx,
		ConnectTimeout: 5 * time.Second,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		ReconnectDelay: DefaultReconnectDelay,
		Logger:         newLogger(slog.Debug),
	}

	c, err := conf.New("tcp", "127.0.0.1:4002")
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}

	go run(c)

	if err := c.Wait(); err != nil {
		log.Fatal(err)
	}
}

func run(c *nanorpc.Client) {
	for {
		select {
		case err, ok := <-c.Pong():
			if !reconnect.IsNoiseError(err) {
				log.Println("pong?", ok, err)
			}
		case <-time.After(1 * time.Second):
			log.Println(".")
		}
	}
}
