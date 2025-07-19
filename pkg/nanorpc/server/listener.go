package server

import "net"

// ListenerAdapter wraps net.Listener to implement our Listener interface
type ListenerAdapter struct {
	net.Listener
}

// NewListenerAdapter creates a new listener adapter
func NewListenerAdapter(listener net.Listener) *ListenerAdapter {
	return &ListenerAdapter{Listener: listener}
}
