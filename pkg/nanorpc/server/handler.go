package server

import (
	"context"

	"github.com/amery/nanorpc/pkg/nanorpc"
)

// DefaultMessageHandler implements MessageHandler interface
type DefaultMessageHandler struct{}

// NewDefaultMessageHandler creates a new message handler
func NewDefaultMessageHandler() *DefaultMessageHandler {
	return &DefaultMessageHandler{}
}

// HandleMessage processes a decoded request
func (h *DefaultMessageHandler) HandleMessage(ctx context.Context, session Session, req *nanorpc.NanoRPCRequest) error {
	switch req.RequestType {
	case nanorpc.NanoRPCRequest_TYPE_PING:
		return h.handlePing(session, req)
	default:
		// Ignore unsupported request types for now
		return nil
	}
}

// handlePing processes ping requests and sends pong responses
func (h *DefaultMessageHandler) handlePing(session Session, req *nanorpc.NanoRPCRequest) error {
	response := &nanorpc.NanoRPCResponse{
		RequestId:      req.RequestId,
		ResponseType:   nanorpc.NanoRPCResponse_TYPE_PONG,
		ResponseStatus: nanorpc.NanoRPCResponse_STATUS_OK,
	}

	responseData, err := nanorpc.EncodeResponse(response, nil)
	if err != nil {
		return err
	}

	// Cast to get Write method
	if sessionWriter, ok := session.(*DefaultSession); ok {
		_, err = sessionWriter.Write(responseData)
		return err
	}

	return nil
}
