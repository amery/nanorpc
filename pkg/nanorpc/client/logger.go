package client

import (
	"net"

	"darvaza.org/slog"
	"darvaza.org/slog/handlers/discard"

	"github.com/amery/nanorpc/pkg/nanorpc/common"
)

// getLogger returns the base logger for the client, creating one if needed
func (c *Client) getLogger() slog.Logger {
	if c.logger == nil {
		// Create a simple discard logger if none provided
		c.logger = discard.New()
	}
	return c.logger
}

// WithDebug returns an annotated debug-level logger
func (c *Client) WithDebug(addr net.Addr) (slog.Logger, bool) {
	logger := c.getLogger()
	if debug, ok := logger.Debug().WithEnabled(); ok {
		return common.WithRemoteAddr(debug, addr), true
	}
	return nil, false
}

// LogDebug writes a log entry at debug-level.
func (c *Client) LogDebug(addr net.Addr, msg string) {
	if l, ok := c.WithDebug(addr); ok {
		l.Print(msg)
	}
}

// WithInfo returns an annotated info-level logger
func (c *Client) WithInfo(addr net.Addr) (slog.Logger, bool) {
	logger := c.getLogger()
	if info, ok := logger.Info().WithEnabled(); ok {
		return common.WithRemoteAddr(info, addr), true
	}
	return nil, false
}

// LogInfo writes a log entry at info-level.
func (c *Client) LogInfo(addr net.Addr, msg string) {
	if l, ok := c.WithInfo(addr); ok {
		l.Print(msg)
	}
}

// WithWarn returns an annotated warn-level logger
func (c *Client) WithWarn(addr net.Addr, err error) (slog.Logger, bool) {
	logger := c.getLogger()
	if warn, ok := logger.Warn().WithEnabled(); ok {
		warn = common.WithError(warn, err)
		return common.WithRemoteAddr(warn, addr), true
	}
	return nil, false
}

// LogWarn writes a log entry at warn-level.
func (c *Client) LogWarn(addr net.Addr, err error, msg string) {
	if l, ok := c.WithWarn(addr, err); ok {
		l.Print(msg)
	}
}

// WithError returns an annotated error-level logger
func (c *Client) WithError(addr net.Addr, err error) (slog.Logger, bool) {
	logger := c.getLogger()
	if errorLog, ok := logger.Error().WithEnabled(); ok {
		errorLog = common.WithError(errorLog, err)
		return common.WithRemoteAddr(errorLog, addr), true
	}
	return nil, false
}

// getErrorLogger returns an error-level logger without address (for internal use)
func (c *Client) getErrorLogger(err error) (slog.Logger, bool) {
	logger := c.getLogger()
	if errorLog, ok := logger.Error().WithEnabled(); ok {
		return common.WithError(errorLog, err), true
	}
	return nil, false
}

// LogError writes a log entry at error-level.
func (c *Client) LogError(addr net.Addr, err error, msg string) {
	if l, ok := c.WithError(addr, err); ok {
		l.Print(msg)
	}
}

// getLogger returns the configured session logger or lazily initializes one
func (cs *Session) getLogger() slog.Logger {
	if cs.logger == nil {
		// Fallback initialization if logger wasn't set during creation
		logger := common.WithComponent(cs.c.getLogger(), common.ComponentSession)
		logger = common.WithRemoteAddr(logger, cs.ra)
		cs.logger = logger
	}
	return cs.logger
}

// WithDebug returns an annotated debug-level logger
func (cs *Session) WithDebug() (slog.Logger, bool) {
	logger := cs.getLogger()
	if debug, ok := logger.Debug().WithEnabled(); ok {
		return debug, true
	}
	return nil, false
}

// LogDebug writes a log entry at debug-level.
func (cs *Session) LogDebug(msg string) {
	if l, ok := cs.WithDebug(); ok {
		l.Print(msg)
	}
}

// WithInfo returns an annotated info-level logger
func (cs *Session) WithInfo() (slog.Logger, bool) {
	logger := cs.getLogger()
	if info, ok := logger.Info().WithEnabled(); ok {
		return info, true
	}
	return nil, false
}

// LogInfo writes a log entry at info-level.
func (cs *Session) LogInfo(msg string) {
	if l, ok := cs.WithInfo(); ok {
		l.Print(msg)
	}
}

// WithWarn returns an annotated warn-level logger
func (cs *Session) WithWarn(err error) (slog.Logger, bool) {
	logger := cs.getLogger()
	if warn, ok := logger.Warn().WithEnabled(); ok {
		return common.WithError(warn, err), true
	}
	return nil, false
}

// LogWarn writes a log entry at warn-level.
func (cs *Session) LogWarn(err error, msg string) {
	if l, ok := cs.WithWarn(err); ok {
		l.Print(msg)
	}
}

// WithError returns an annotated error-level logger
func (cs *Session) WithError(err error) (slog.Logger, bool) {
	logger := cs.getLogger()
	if errorLog, ok := logger.Error().WithEnabled(); ok {
		return common.WithError(errorLog, err), true
	}
	return nil, false
}

// LogError writes a log entry at error-level.
func (cs *Session) LogError(err error, msg string) {
	if l, ok := cs.WithError(err); ok {
		l.Print(msg)
	}
}
