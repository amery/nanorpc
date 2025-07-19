package server

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/amery/nanorpc/pkg/nanorpc"
)

func TestDecoupledServer_PingPong(t *testing.T) {
	// Create a test listener
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close()

	// Create decoupled server
	server := NewDefaultServer(listener)

	// Start server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.Serve(ctx)
	}()

	// Give server time to start
	time.Sleep(50 * time.Millisecond)

	// Connect to server
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Send ping request
	pingReq := &nanorpc.NanoRPCRequest{
		RequestId:   789,
		RequestType: nanorpc.NanoRPCRequest_TYPE_PING,
	}

	pingData, err := nanorpc.EncodeRequest(pingReq, nil)
	if err != nil {
		t.Fatalf("Failed to encode ping: %v", err)
	}

	if _, err := conn.Write(pingData); err != nil {
		t.Fatalf("Failed to send ping: %v", err)
	}

	// Read pong response
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	// Decode and verify response
	response, _, err := nanorpc.DecodeResponse(buffer[:n])
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ResponseType != nanorpc.NanoRPCResponse_TYPE_PONG {
		t.Fatalf("Expected PONG, got %v", response.ResponseType)
	}

	if response.RequestId != 789 {
		t.Fatalf("Expected RequestId=789, got %d", response.RequestId)
	}

	if response.ResponseStatus != nanorpc.NanoRPCResponse_STATUS_OK {
		t.Fatalf("Expected STATUS_OK, got %v", response.ResponseStatus)
	}

	// Shutdown server
	cancel()

	select {
	case err := <-serverErr:
		if err != nil && err != context.Canceled {
			t.Fatalf("Server stopped with unexpected error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Server shutdown timeout")
	}
}

func TestDecoupledServer_Shutdown(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	server := NewDefaultServer(listener)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		t.Fatalf("Failed to shutdown server: %v", err)
	}
}
