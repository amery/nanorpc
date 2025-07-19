package server

import (
	"bufio"
	"context"
	"fmt"
	"net"

	"github.com/amery/nanorpc/pkg/nanorpc"
)

// DefaultSession implements Session interface
type DefaultSession struct {
	id         string
	conn       net.Conn
	handler    MessageHandler
	remoteAddr string
}

// NewDefaultSession creates a new session
func NewDefaultSession(conn net.Conn, handler MessageHandler) *DefaultSession {
	return &DefaultSession{
		id:         generateSessionID(conn),
		conn:       conn,
		handler:    handler,
		remoteAddr: conn.RemoteAddr().String(),
	}
}

// ID returns the session identifier
func (s *DefaultSession) ID() string {
	return s.id
}

// RemoteAddr returns the remote address
func (s *DefaultSession) RemoteAddr() string {
	return s.remoteAddr
}

// Handle processes messages for this session
func (s *DefaultSession) Handle(ctx context.Context) error {
	defer s.Close()

	scanner := bufio.NewScanner(s.conn)
	scanner.Split(nanorpc.Split)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("scan error: %w", err)
			}
			return nil // EOF
		}

		req, _, err := nanorpc.DecodeRequest(scanner.Bytes())
		if err != nil {
			// Log decode error but continue
			continue
		}

		if err := s.handler.HandleMessage(ctx, s, req); err != nil {
			// Log handler error but continue
			continue
		}
	}
}

// Close closes the session
func (s *DefaultSession) Close() error {
	return s.conn.Close()
}

// Write sends data to the client
func (s *DefaultSession) Write(data []byte) (int, error) {
	return s.conn.Write(data)
}

// generateSessionID creates a unique session identifier
func generateSessionID(conn net.Conn) string {
	return fmt.Sprintf("session-%s", conn.RemoteAddr().String())
}
