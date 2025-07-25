package server

import (
	"context"
	"sync"

	"darvaza.org/core"

	"github.com/amery/nanorpc/pkg/nanorpc"
)

// RequestContext provides request information and response utilities
type RequestContext struct {
	Session  Session
	Request  *nanorpc.NanoRPCRequest
	Path     string // Resolved path (from string or hash)
	PathHash uint32 // The hash of the path (computed or provided)
}

// DefaultMessageHandler implements MessageHandler interface
type DefaultMessageHandler struct {
	handlers  map[string]RequestHandler
	hashCache *nanorpc.HashCache
	mu        sync.RWMutex
}

// NewDefaultMessageHandler creates a new message handler with an optional HashCache.
// If hashCache is nil, a new one will be created.
func NewDefaultMessageHandler(hashCache *nanorpc.HashCache) *DefaultMessageHandler {
	if hashCache == nil {
		hashCache = &nanorpc.HashCache{}
	}
	return &DefaultMessageHandler{
		handlers:  make(map[string]RequestHandler),
		hashCache: hashCache,
	}
}

// RegisterHandlerFunc registers a handler function for a specific path
func (h *DefaultMessageHandler) RegisterHandlerFunc(path string, fn RequestHandlerFunc) error {
	return h.RegisterHandler(path, fn)
}

// RegisterHandler registers a handler for a specific path
func (h *DefaultMessageHandler) RegisterHandler(path string, handler RequestHandler) error {
	if h == nil {
		return core.ErrNilReceiver
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.handlers == nil {
		h.handlers = make(map[string]RequestHandler)
	}

	if handler == nil {
		return h.doUnregister(path)
	}

	return h.doRegister(path, handler)
}

func (h *DefaultMessageHandler) doUnregister(path string) error {
	if _, exists := h.handlers[path]; exists {
		delete(h.handlers, path)
		return nil
	}
	return core.ErrNotExists
}

func (h *DefaultMessageHandler) doRegister(path string, handler RequestHandler) error {
	if _, exists := h.handlers[path]; exists {
		return core.ErrExists
	}

	// Populate the hash cache with this path
	if _, err := h.hashCache.Hash(path); err != nil {
		return err
	}

	h.handlers[path] = handler
	return nil
}

// HandleMessage processes a decoded request
func (h *DefaultMessageHandler) HandleMessage(ctx context.Context, session Session, req *nanorpc.NanoRPCRequest) error {
	switch req.RequestType {
	case nanorpc.NanoRPCRequest_TYPE_PING:
		return h.handlePing(ctx, session, req)
	case nanorpc.NanoRPCRequest_TYPE_REQUEST:
		return h.handleRequest(ctx, session, req)
	default:
		// Ignore unsupported request types for now
		return nil
	}
}

// sendErrorResponse is a helper to send an error response
func sendErrorResponse(session Session, req *nanorpc.NanoRPCRequest,
	status nanorpc.NanoRPCResponse_Status, message string) error {
	response := &nanorpc.NanoRPCResponse{
		RequestId:       req.RequestId,
		ResponseType:    nanorpc.NanoRPCResponse_TYPE_RESPONSE,
		ResponseStatus:  status,
		ResponseMessage: message,
	}
	return session.SendResponse(req, response)
}

// handlePing processes ping requests and sends pong responses
func (*DefaultMessageHandler) handlePing(_ context.Context, session Session, req *nanorpc.NanoRPCRequest) error {
	response := &nanorpc.NanoRPCResponse{
		RequestId:      req.RequestId,
		ResponseType:   nanorpc.NanoRPCResponse_TYPE_PONG,
		ResponseStatus: nanorpc.NanoRPCResponse_STATUS_OK,
	}

	return session.SendResponse(req, response)
}

// handleRequest processes TYPE_REQUEST messages
func (h *DefaultMessageHandler) handleRequest(ctx context.Context, session Session, req *nanorpc.NanoRPCRequest) error {
	// Extract path and hash from request
	var path string
	var pathHash uint32

	switch p := req.PathOneof.(type) {
	case *nanorpc.NanoRPCRequest_Path:
		path = p.Path
		var err error
		pathHash, err = h.hashCache.Hash(path) // Compute and cache the hash
		if err != nil {
			// Hash collision on incoming request - return internal error
			return sendErrorResponse(session, req,
				nanorpc.NanoRPCResponse_STATUS_INTERNAL_ERROR,
				"path hash collision")
		}
	case *nanorpc.NanoRPCRequest_PathHash:
		pathHash = p.PathHash
		// Try to resolve hash to path
		if resolvedPath, ok := h.hashCache.Path(pathHash); ok {
			path = resolvedPath
		}
		// If we can't resolve the hash, path remains empty
	}

	// Look up handler
	h.mu.RLock()
	handler, exists := h.handlers[path]
	h.mu.RUnlock()

	if !exists || handler == nil || path == "" {
		// No handler registered or path couldn't be resolved
		return sendErrorResponse(session, req,
			nanorpc.NanoRPCResponse_STATUS_NOT_FOUND,
			"no handler registered for path")
	}

	// Create request context
	reqCtx := &RequestContext{
		Session:  session,
		Request:  req,
		Path:     path,
		PathHash: pathHash,
	}

	// Call the handler
	return handler.Handle(ctx, reqCtx)
}
