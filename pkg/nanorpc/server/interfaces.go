package server

import (
	"context"
	"net"

	"github.com/amery/nanorpc/pkg/nanorpc"
)

// Listener handles connection acceptance
type Listener interface {
	// Accept waits for and returns the next connection
	Accept() (net.Conn, error)
	// Close closes the listener
	Close() error
	// Addr returns the listener's network address
	Addr() net.Addr
}

// SessionManager manages connection lifecycle and tracking
type SessionManager interface {
	// AddSession creates a new session for the connection
	AddSession(conn net.Conn) Session
	// RemoveSession removes a session by ID
	RemoveSession(sessionID string)
	// GetSession retrieves a session by ID
	GetSession(sessionID string) Session
	// Shutdown gracefully closes all sessions
	Shutdown(ctx context.Context) error
}

// Session represents a single client connection
type Session interface {
	// ID returns the unique session identifier
	ID() string
	// RemoteAddr returns the remote address
	RemoteAddr() string
	// Handle processes messages for this session
	Handle(ctx context.Context) error
	// Close closes the session
	Close() error
}

// MessageHandler processes protocol messages
type MessageHandler interface {
	// HandleMessage processes a decoded request
	HandleMessage(ctx context.Context, session Session, req *nanorpc.NanoRPCRequest) error
}
