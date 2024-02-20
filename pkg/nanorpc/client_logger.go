package nanorpc

import (
	"net"

	"darvaza.org/slog"
)

func (c *Client) info(addr net.Addr) (slog.Logger, bool) {
	if log := c.options.Logger; log != nil {
		if l, ok := log.Info().WithEnabled(); ok {
			if addr != nil {
				s := addr.String()
				if s != "" {
					l = l.WithField("addr", s)
				}
			}

			return l, true
		}
	}

	return nil, false
}

func (c *Client) error(addr net.Addr, err error) (slog.Logger, bool) {
	if log := c.options.Logger; log != nil {
		if l, ok := log.Error().WithEnabled(); ok {
			if addr != nil {
				s := addr.String()
				if s != "" {
					l = l.WithField("addr", s)
				}
			}

			if err != nil {
				l = l.WithField(slog.ErrorFieldName, err)
			}

			return l, true
		}
	}

	return nil, false
}

func (c *Client) say(conn net.Conn, format string, args ...any) {
	var ra net.Addr
	if conn != nil {
		ra = conn.RemoteAddr()
	}

	if l, ok := c.info(ra); ok {
		l.Printf(format, args...)
	}
}

func (c *Client) sayError(conn net.Conn, err error) {
	var ra net.Addr
	if conn != nil {
		ra = conn.RemoteAddr()
	}

	if l, ok := c.error(ra, err); ok {
		l.Print()
	}
}
