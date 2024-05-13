package nanorpc

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"darvaza.org/core"
)

var (
	// ErrNoCallback tells a callback function is expected but wasn't provided
	ErrNoCallback = core.Wrap(os.ErrInvalid, "missing callback")

	// ErrNotConnected tells the [Client] isn't currently connected
	ErrNotConnected = core.Wrap(os.ErrClosed, "not connected")

	// ErrNoResponse indicates the server didn't answer before disconnection
	ErrNoResponse = core.NewTimeoutError(errors.New("no response"))
)

// AsError extracts an error from the
// status of a response.
func AsError(res *NanoRPCResponse) error {
	if res != nil {
		switch res.ResponseStatus {
		case NanoRPCResponse_STATUS_OK:
			return nil
		case NanoRPCResponse_STATUS_NOT_FOUND:
			return fs.ErrNotExist
		case NanoRPCResponse_STATUS_NOT_AUTHORIZED:
			return fs.ErrPermission
		default:
			return fmt.Errorf("invalid state %v", int(res.ResponseStatus))
		}
	}
	return ErrNoResponse
}
